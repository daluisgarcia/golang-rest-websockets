package server

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/daluisgarcia/golang-rest-websockets/database"
	"github.com/daluisgarcia/golang-rest-websockets/repositories"
	"github.com/daluisgarcia/golang-rest-websockets/websockets"
	"github.com/gorilla/mux"
)

type Config struct {
	Port        string
	JWTSecret   string
	DatabaseUrl string
}

type Server interface {
	Config() *Config
	Hub() *websockets.Hub
}

type Broker struct {
	config *Config
	router *mux.Router
	hub    *websockets.Hub
}

func NewServer(ctx context.Context, config *Config) (*Broker, error) {
	if config.Port == "" {
		return nil, fmt.Errorf("port is required")
	}

	if config.JWTSecret == "" {
		return nil, fmt.Errorf("jwt secret is required")
	}

	if config.DatabaseUrl == "" {
		return nil, fmt.Errorf("database url is required")
	}

	return &Broker{
		config: config,
		router: mux.NewRouter(),
		hub:    websockets.NewHub(),
	}, nil
}

func (b *Broker) Config() *Config {
	return b.config
}

func (b *Broker) Hub() *websockets.Hub {
	return b.hub
}

func (b *Broker) Start(binder func(s Server, r *mux.Router)) {
	b.router = mux.NewRouter()
	b.hub = websockets.NewHub()
	binder(b, b.router)
	repo, err := database.NewPostgresRepository(b.config.DatabaseUrl)

	if err != nil {
		log.Fatal(err)
	}

	repositories.SetRepository(repo)

	log.Println("Server started on port", b.config.Port)

	if err := http.ListenAndServe(":"+b.config.Port, b.router); err != nil {
		log.Fatal("Error when starting the server: ", err)
	}

}
