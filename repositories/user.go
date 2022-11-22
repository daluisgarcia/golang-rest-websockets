package repositories

import (
	"context"

	"github.com/daluisgarcia/golang-rest-websockets/models"
)

type UserRepository interface {
	InsertUser(ctx context.Context, user *models.User) error
	FindUserById(ctx context.Context, id string) (*models.User, error)
	Close() error
}

var implementation UserRepository

func SetRepository(repo UserRepository) {
	implementation = repo
}

func InsertUser(ctx context.Context, user *models.User) error {
	return implementation.InsertUser(ctx, user)
}

func FindUserById(ctx context.Context, id string) (*models.User, error) {
	return implementation.FindUserById(ctx, id)
}

func Close() error {
	return implementation.Close()
}
