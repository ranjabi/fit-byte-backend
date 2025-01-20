package main

import (
	"context"
	"fit-byte/constants"
	"fit-byte/db"
	"fit-byte/usecases/activity"
	"fit-byte/usecases/auth"
	"fit-byte/usecases/user"
	"fit-byte/utils"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/jwtauth/v5"
	"github.com/joho/godotenv"
)

func main() {
    err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

    ctx := context.Background()
	pgConn := db.Setup(ctx)

    userRepository := user.NewUserRepository(ctx, pgConn)
	activityRepository := activity.NewActivityRepository(ctx, pgConn)

    authService := auth.NewAuthService(userRepository)
    userService := user.NewUserService(userRepository)
	activityService := activity.NewActivityService(activityRepository)
    
    authHandler := auth.NewAuthHandler(authService)
    userHandler := user.NewUserHandler(userService)
	activityHandler := activity.NewActivityHandler(activityService)
    
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
		})
	})
    
    http.ListenAndServe(":8080", r)
}
