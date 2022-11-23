package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/daluisgarcia/golang-rest-websockets/handlers"
	"github.com/daluisgarcia/golang-rest-websockets/middleware"
	"github.com/daluisgarcia/golang-rest-websockets/server"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

func BindRoutes(s server.Server, r *mux.Router) {
	// Free of authentication routes
	r.HandleFunc("/", handlers.HomeHandler(s)).Methods("GET")
	r.HandleFunc("/signup", handlers.SignUpHandler(s)).Methods(http.MethodPost)
	r.HandleFunc("/login", handlers.LoginHandler(s)).Methods(http.MethodPost)
	r.HandleFunc("/posts/{id}", handlers.GetPostHandler(s)).Methods(http.MethodGet)
	r.HandleFunc("/posts", handlers.ListPostsHandler(s)).Methods(http.MethodGet)
	r.HandleFunc("/ws", s.Hub().HandleWebSocket)

	api := r.PathPrefix("/api/v1").Subrouter() // Defining a subrouter for the API

	api.Use(middleware.CheckAuthMiddleware(s)) // Applies a middleware to all routes of the api

	api.HandleFunc("/me", handlers.MeHandler(s)).Methods(http.MethodGet)
	api.HandleFunc("/posts", handlers.InsertPostHandler(s)).Methods(http.MethodPost)
	api.HandleFunc("/posts/{id}", handlers.UpdatePostHandler(s)).Methods(http.MethodPut)
	api.HandleFunc("/posts/{id}", handlers.DeletePostHandler(s)).Methods(http.MethodDelete)
}

func main() {
	err := godotenv.Load()

	if err != nil {
		log.Fatal("Error loading .env file")
	}

	PORT := os.Getenv("PORT")
	JWT_SECRET := os.Getenv("JWT_SECRET")
	DATABASE_URL := os.Getenv("DATABASE_URL")

	s, err := server.NewServer(context.Background(), &server.Config{
		Port:        PORT,
		JWTSecret:   JWT_SECRET,
		DatabaseUrl: DATABASE_URL,
	})

	if err != nil {
		log.Fatal(err)
	}

	s.Start(BindRoutes)

}
