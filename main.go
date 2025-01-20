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

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/jwtauth/v5"
	"github.com/joho/godotenv"
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

func main() {
    err := godotenv.Load(".env")
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
