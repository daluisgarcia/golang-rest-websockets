package server

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type Config struct {
	Port        string
	JWTSecret   string
	DatabaseUrl string
}

type Server interface {
	Config() *Config
}

type Broker struct {
	config *Config
	router *mux.Router
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
	}, nil
}

func (b *Broker) Config() *Config {
	return b.config
}

func (b *Broker) Start(binder func(s Server, r *mux.Router)) {
	b.router = mux.NewRouter()
	binder(b, b.router)
	log.Println("Server started on port", b.config.Port)

	if err := http.ListenAndServe(":"+b.config.Port, b.router); err != nil {
		log.Fatal("Error when starting the server: ", err)
	}

}
