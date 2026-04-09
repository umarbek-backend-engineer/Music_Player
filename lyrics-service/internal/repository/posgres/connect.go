package posgres

import (
	"context"
	"fmt"
	"lyrics-service/internal/config"
	"time"

	"github.com/jackc/pgx/v5"
)

func Connect() (*pgx.Conn, error) {
	cgf := config.Load()

	posgresURL := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=diable", cgf.DB_User, cgf.DB_Password, cgf.Api_Host, cgf.DB_Port, cgf.DB_Name)
	client, err := pgx.Connect(context.Background(), posgresURL)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	if err := client.Ping(ctx); err != nil {
		return nil, err
	}

	return client, nil
}
