package db_connect

import (
	"context"
	"fmt"
	"music-service/internal/config"
	"time"

	"github.com/jackc/pgx/v5"
)

// creates posgres client
func Connect() (*pgx.Conn, error) {
	cgf := config.Load()
	postgresURL := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", cgf.DB_User, cgf.DB_Password, cgf.DB_Host, cgf.DB_Port, cgf.DB_Name)
	client, err := pgx.Connect(context.Background(), postgresURL)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	// checking if the client can connect to the posgres server.
	if err := client.Ping(ctx); err != nil {
		return nil, err
	}

	return client, nil
}
