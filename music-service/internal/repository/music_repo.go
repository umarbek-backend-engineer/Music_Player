package repository

import (
	"context"
	"music-service/internal/repository/db_connect.go"

	"github.com/google/uuid"
)

func UploadMusicDBHandler(ctx context.Context, filename, filepath string) (uuid.UUID, error) {

	var id uuid.UUID

	conn, err := db_connect.Connect()
	if err != nil {
		return uuid.Nil, err
	}
	defer conn.Close(ctx)

	err = conn.QueryRow(ctx, "insert into musics (filename, filepath) values ($1, $2) returning id", filename, filepath).Scan(&id)
	if err != nil {
		return uuid.Nil, err
	}

	return id, nil
}
