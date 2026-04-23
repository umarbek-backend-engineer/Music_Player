package postgres

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/umarbek-backend-engineer/Music_Player/internal/config"
)

func Connect() (*pgx.Conn, error) {

	cgf := config.Load()

	// connect to postgres 
	postgreURL := fmt.Sprintf("postgresql://%s:%s@%s:%s/%s?sslmode=disable", cgf.DB_user, cgf.DB_password, cgf.DB_host, cgf.DB_port, cgf.DB_name)
	client, err := pgx.Connect(context.Background(), postgreURL)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	// check the conection
	if err := client.Ping(ctx); err != nil {
		return nil, err
	}

	return client, nil
}
