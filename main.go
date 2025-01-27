package main

import (
	"context"
	"fit-byte/constants"
	"fit-byte/db"
	"fit-byte/usecases/activity"
	"fit-byte/usecases/auth"
	"fit-byte/usecases/file"
	"fit-byte/usecases/user"
	"fit-byte/utils"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/jwtauth/v5"
	"github.com/joho/godotenv"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var s3Client *s3.Client

func initS3(ctx context.Context) error {
	awsConfig, err := config.LoadDefaultConfig(ctx, config.WithRegion(os.Getenv("AWS_REGION")))
	if err != nil {
		log.Fatalf("Unable to load AWS config: %v", err)
	}

	s3Client = s3.NewFromConfig(awsConfig)
	return nil
}

type Metrics struct {
	httpRequestsTotal *prometheus.CounterVec
	httpDuration *prometheus.HistogramVec
}

func NewMetrics(reg prometheus.Registerer) *Metrics {
	m := &Metrics{
		httpRequestsTotal: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: "http_requests_total",
				Help: "Total number of HTTP requests.",
			},
			[]string{"method", "status", "handler"},
		),
		httpDuration: prometheus.NewHistogramVec(
			prometheus.HistogramOpts{
				Name:    "http_request_duration_seconds",
				Help:    "Histogram of HTTP request durations.",
				Buckets: prometheus.DefBuckets,
			},
			[]string{"method", "handler"},
		),
	}
	reg.MustRegister(m.httpRequestsTotal)
	reg.MustRegister(m.httpDuration)
	return m
}

type AppResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

func NewResponseWriter(w http.ResponseWriter) *AppResponseWriter {
	return &AppResponseWriter{w, http.StatusOK}
}

func (rw *AppResponseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

func (m *Metrics) MetricsMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        start := time.Now()

		rw := NewResponseWriter(w)
        next.ServeHTTP(rw, r)

        duration := time.Since(start).Seconds()

        // Automatically set the labels based on the request
        m.httpRequestsTotal.WithLabelValues(r.Method, http.StatusText(rw.statusCode), r.URL.Path).Inc()
        m.httpDuration.WithLabelValues(r.Method, r.URL.Path).Observe(duration)
    })
}

func main() {
	env := ".env"
	if os.Getenv("ENV") == "dev.docker" {
		env = ".env.dev.docker"
	}

    err := godotenv.Load(env)
	if err != nil {
		log.Fatal("Error loading .env file")
	}

    ctx := context.Background()
	pgConn := db.Setup(ctx)
	if err := initS3(ctx); err != nil {
		log.Fatal(err.Error())
	}

    userRepository := user.NewUserRepository(ctx, pgConn)
	activityRepository := activity.NewActivityRepository(ctx, pgConn)

    authService := auth.NewAuthService(userRepository)
    userService := user.NewUserService(userRepository)
	activityService := activity.NewActivityService(activityRepository)
	fileService := file.NewFileService(s3Client, ctx)
    
    authHandler := auth.NewAuthHandler(authService)
    userHandler := user.NewUserHandler(userService)
	activityHandler := activity.NewActivityHandler(activityService)
	fileHandler := file.NewFileHandler(fileService)
    
    r := chi.NewRouter()
    r.Use(middleware.Logger)
	r.Use(middleware.Heartbeat("/ping"))
	
	m := NewMetrics(prometheus.DefaultRegisterer)
	r.Use(m.MetricsMiddleware)
	r.Handle("/metrics", promhttp.Handler())
    
    r.Route("/v1", func(r chi.Router) {
		// public
		r.Group(func(r chi.Router) {
			r.Post("/register", utils.AppHandler(authHandler.HandleRegister))
			r.Post("/login", utils.AppHandler(authHandler.HandleLogin))
		})

		// protected
		r.Group(func(r chi.Router) {
			tokenAuth := jwtauth.New(constants.HASH_ALG, []byte(constants.JWT_SECRET), nil)
			r.Use(jwtauth.Verifier(tokenAuth))
			r.Use(jwtauth.Authenticator(tokenAuth))
			r.Use(utils.AllowContentType("application/json", "multipart/form-data"))

			r.Get("/user", utils.AppHandler(userHandler.HandleGetUser))
			r.Patch("/user", utils.AppHandler(userHandler.HandleUpdateUser))

			r.Get("/activity", utils.AppHandler(activityHandler.HandleGetAllActivities))
			r.Post("/activity", utils.AppHandler(activityHandler.HandleCreateActivity))
			r.Patch("/activity/{activityId}", utils.AppHandler(activityHandler.HandleUpdateActivity))
			r.Delete("/activity/{activityId}", utils.AppHandler(activityHandler.HandleDeleteActivity))

			r.Post("/file", utils.AppHandler(fileHandler.HandleUploadFile))
		})
	})
    
    http.ListenAndServe(":8080", r)
}
