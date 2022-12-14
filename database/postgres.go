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

func (repo *PostgresRepository) Close() error {
	return repo.db.Close()
}

func (repo *PostgresRepository) InsertUser(ctx context.Context, user *models.User) error {
	_, err := repo.db.ExecContext(ctx, "INSERT INTO users (id, email, password) VALUES ($1, $2, $3)", user.Id, user.Email, user.Password)
	return err
}

func (repo *PostgresRepository) FindUserById(ctx context.Context, id string) (*models.User, error) {
	rows, err := repo.db.QueryContext(ctx, "SELECT id, email FROM users WHERE id = $1", id)

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
		if err := rows.Scan(&user.Id, &user.Email); err == nil {
			return &user, err
		}
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return &user, nil
}

func (repo *PostgresRepository) FindUserByEmail(ctx context.Context, email string) (*models.User, error) {
	rows, err := repo.db.QueryContext(ctx, "SELECT id, email, password FROM users WHERE email = $1", email)

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

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return &user, nil
}

func (repo *PostgresRepository) InsertPost(ctx context.Context, post *models.Post) error {
	_, err := repo.db.ExecContext(ctx, "INSERT INTO posts (id, post_content, user_id) VALUES ($1, $2, $3)", post.Id, post.PostContent, post.UserId)
	return err
}

func (repo *PostgresRepository) FindPostById(ctx context.Context, id string) (*models.Post, error) {
	rows, err := repo.db.QueryContext(ctx, "SELECT id, post_content, user_id, created_at FROM posts WHERE id = $1", id)

	defer func() { // Alows to validate the error after the function returns
		err := rows.Close()

		if err != nil {
			log.Fatal(err)
		}
	}()

	if err != nil {
		return nil, err
	}

	var post = models.Post{}
	for rows.Next() {
		if err := rows.Scan(&post.Id, &post.PostContent, &post.UserId, &post.CreatedAt); err == nil {
			return &post, err
		}
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return &post, nil
}

func (repo *PostgresRepository) UpdatePost(ctx context.Context, post *models.Post) error {
	_, err := repo.db.ExecContext(ctx, "UPDATE posts SET post_content = $1 WHERE id = $2 AND user_id = $3", post.PostContent, post.Id, post.UserId)
	return err
}

func (repo *PostgresRepository) DeletePost(ctx context.Context, id string, userId string) error {
	_, err := repo.db.ExecContext(ctx, "DELETE FROM posts WHERE id = $1 AND user_id = $2", id, userId)
	return err
}

func (repo *PostgresRepository) ListPosts(ctx context.Context, page uint64, userId string) ([]*models.Post, error) {
	var postPerPage int64 = 10
	if page <= 0 {
		page = 1
	}

	rows, err := repo.db.QueryContext(
		ctx,
		"SELECT id, post_content, user_id, created_at FROM posts WHERE user_id = $1 ORDER BY created_at DESC LIMIT $2 OFFSET $3",
		userId, postPerPage, (page-1)*uint64(postPerPage),
	)

	defer func() { // Alows to validate the error after the function returns
		err := rows.Close()

		if err != nil {
			log.Fatal(err)
		}
	}()

	if err != nil {
		return nil, err
	}

	var posts []*models.Post
	for rows.Next() {
		var post = models.Post{}
		if err := rows.Scan(&post.Id, &post.PostContent, &post.UserId, &post.CreatedAt); err == nil {
			posts = append(posts, &post)
		}
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return posts, nil
}
