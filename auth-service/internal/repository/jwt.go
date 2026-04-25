package repository

import (
	"context"
	"time"

	"github.com/umarbek-backend-engineer/Music_Player/internal/config"
	"github.com/umarbek-backend-engineer/Music_Player/internal/repository/postgres"
)

// a crud operation which will save the refresh token inside the database and set exiration date
func InsertRefreshToken(ctx context.Context, id, token, user_agent, ip_address string) error {
	// load configurations
	cgf := config.Load()
	// connect to database
	conn, err := postgres.Connect()
	if err != nil {
		return err
	}
	defer conn.Close(ctx)

	// Get the duration for refresh token
	duration, err := time.ParseDuration(cgf.REF_JWT_exp)
	if err != nil {
		return err
	}

	// set the time with with duration: now + duration
	expires_at := time.Now().Add(duration)

	// giving the query to database to save the refresh token in database
	_, err = conn.Exec(ctx, "insert into sessions (user_id, token_hash, user_agent, ip_address, expires_at) values ($1,$2,$3)", id, token, user_agent, ip_address, expires_at)
	if err != nil {
		return err
	}

	return nil
}

// a crud operation which will save the refresh token inside the database and set exiration date
func UpdateRefreshToken(ctx context.Context, id, token string) error {
	// load configurations
	cgf := config.Load()
	// connect to database
	conn, err := postgres.Connect()
	if err != nil {
		return err
	}
	defer conn.Close(ctx)

	// Get the duration for refresh token
	duration, err := time.ParseDuration(cgf.REF_JWT_exp)
	if err != nil {
		return err
	}

	// set the time with with duration: now + duration
	expires_at := time.Now().Add(duration)

	// giving the query to database to save the refresh token in database
	_, err = conn.Exec(ctx, "update tabel sessions set token_hash = $1, expires_at = $2 where user_id = $3", id, token, expires_at)
	if err != nil {
		return err
	}

	// return nil if there is not problem in updating
	return nil
}
