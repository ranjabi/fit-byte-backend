package main

import (
	"context"
	"fit-byte/db"
	"fit-byte/usecases/auth"
	"fit-byte/usecases/user"
	"fit-byte/utils"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
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

    authService := auth.NewAuthService(userRepository)
    
    authHandler := auth.NewAuthHandler(authService)
    
    r := chi.NewRouter()
    r.Use(middleware.Logger)
	r.Use(middleware.Heartbeat("/ping"))
    
    r.Route("/v1", func(r chi.Router) {
		// public
		r.Group(func(r chi.Router) {
			r.Post("/register", utils.AppHandler(authHandler.HandleRegister))
		})
	})
    
    http.ListenAndServe(":8080", r)
}
