package repositories

import (
	"context"

	"github.com/daluisgarcia/golang-rest-websockets/models"
)

type Repository interface {
	Close() error
	InsertUser(ctx context.Context, user *models.User) error
	FindUserById(ctx context.Context, id string) (*models.User, error)
	FindUserByEmail(ctx context.Context, email string) (*models.User, error)
	InsertPost(ctx context.Context, post *models.Post) error
	FindPostById(ctx context.Context, id string) (*models.Post, error)
}

var implementation Repository

func SetRepository(repo Repository) {
	implementation = repo
}

func Close() error {
	return implementation.Close()
}

func InsertUser(ctx context.Context, user *models.User) error {
	return implementation.InsertUser(ctx, user)
}

func FindUserById(ctx context.Context, id string) (*models.User, error) {
	return implementation.FindUserById(ctx, id)
}

func FindUserByEmail(ctx context.Context, email string) (*models.User, error) {
	return implementation.FindUserByEmail(ctx, email)
}

func InsertPost(ctx context.Context, post *models.Post) error {
	return implementation.InsertPost(ctx, post)
}

func FindPostById(ctx context.Context, id string) (*models.Post, error) {
	return implementation.FindPostById(ctx, id)
}
