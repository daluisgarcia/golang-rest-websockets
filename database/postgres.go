package database

import (
	"context"
	"database/sql"
	"log"

	"github.com/daluisgarcia/golang-rest-websockets/models"
	_ "github.com/lib/pq"
)

type PostgresRepository struct {
	db *sql.DB
}

func NewPostgresRepository(url string) (*PostgresRepository, error) {
	db, err := sql.Open("postgres", url)

	if err != nil {
		return nil, err
	}

	return &PostgresRepository{
		db: db,
	}, nil
}

func (repo *PostgresRepository) InsertUser(ctx context.Context, user *models.User) error {
	_, err := repo.db.ExecContext(ctx, "INSERT INTO users (email, password) VALUES ($1, $2)", user.Email, user.Password)
	return err
}

func (repo *PostgresRepository) FindUserById(ctx context.Context, id int64) (*models.User, error) {
	rows, err := repo.db.QueryContext(ctx, "SELECT id, email, password FROM users WHERE id = $1", id)

	defer func() { // Alows to validate the error after the function returns
		err := rows.Close()

		if err != nil {
			log.Fatal(err)
		}
	}()

	if err != nil {
		return nil, err
	}

	var user = models.User{}
	for rows.Next() {
		if err := rows.Scan(&user.Id, &user.Email, &user.Password); err == nil {
			return &user, err
		}
	}

	if err == rows.Err() {
		return nil, err
	}

	return &user, nil
}

func (repo *PostgresRepository) Close() error {
	return repo.db.Close()
}
