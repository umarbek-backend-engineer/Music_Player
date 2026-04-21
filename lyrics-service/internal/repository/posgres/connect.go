package posgres

import (
	"context"
	"fmt"

	"time"

	"github.com/jackc/pgx/v5"
	"github.com/umarbek-backend-engineer/Music_Player/lyrics-service/internal/config"
)

func Connect() (*pgx.Conn, error) {
	cgf := config.Load()

	posgresURL := fmt.Sprintf("postgresql://%s:%s@%s:%s/%s?sslmode=disable", cgf.DB_User, cgf.DB_Password, cgf.DB_Host, cgf.DB_Port, cgf.DB_Name)
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
