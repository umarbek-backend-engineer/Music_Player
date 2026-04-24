package repository

import (
	"context"
	"time"

	"github.com/umarbek-backend-engineer/Music_Player/internal/repository/postgres"
)

// a crud operation which will save the refresh token inside the database and set exiration date
func InsertRefreshToken(ctx context.Context, id, token string) error {
	// connect to database
	conn, err := postgres.Connect()
	if err != nil {
		return err
	}
	defer conn.Close(ctx)

	// makein a menual expiration time for refresh token
	expires_at := time.Now().Add(time.Hour * 168)

	// giving the query to database to save the refresh token in database
	_, err = conn.Exec(ctx, "insert into refresh_token (user_id, token_hash, expires_at) values ($1,$2,$3)", id, token, expires_at)
	if err != nil {
		return err
	}

	return nil
}

// a crud operation which will save the refresh token inside the database and set exiration date
func UpdateRefreshToken(ctx context.Context, id, token string) error {
	// connect to database
	conn, err := postgres.Connect()
	if err != nil {
		return err
	}
	defer conn.Close(ctx)

	// makein a menual expiration time for refresh token
	expires_at := time.Now().Add(time.Hour * 168)

	// giving the query to database to save the refresh token in database
	_, err = conn.Exec(ctx, "update tabel refresh_token set token_hash = $1, expires_at = $2 where user_id = $3", id, token, expires_at)
	if err != nil {
		return err
	}

	return nil
}
