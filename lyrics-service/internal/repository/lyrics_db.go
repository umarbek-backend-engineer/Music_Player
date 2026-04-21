package repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/umarbek-backend-engineer/Music_Player/lyrics-service/internal/model"
	"github.com/umarbek-backend-engineer/Music_Player/lyrics-service/internal/repository/posgres"
)

// in this function I will check if the lyrics of the music exists If yes, it will return true, if no, it will return false

func Is_music_lyric_exists(ctx context.Context, music_name string) (bool, error) {
	conn, err := posgres.Connect()
	if err != nil {
		return false, err
	}
	defer conn.Close(ctx)

	var exists bool

	// here i am checking if the same name exists in the data base if yes it will return true, else false
	err = conn.QueryRow(ctx, "select exists (select 1 from lyrics where name = $1)", music_name).Scan(&exists)
	if err != nil && err != sql.ErrNoRows {
		return false, err
	}

	return exists, nil
}

func SaveLyrics(ctx context.Context, musicID, musicName string, content model.Respond) error {
	conn, err := posgres.Connect()
	if err != nil {
		return err
	}
	defer conn.Close(ctx)

	_, err = conn.Exec(ctx,
		"insert into lyrics (music_id, name, content) values ($1, $2, $3)",
		musicID,
		musicName,
		content,
	)

	return err
}

func GetLyricsByMusicID(ctx context.Context, musicID string) (model.Respond, error) {
	conn, err := posgres.Connect()
	if err != nil {
		return model.Respond{}, err
	}
	defer conn.Close(ctx)

	var content model.Respond
	err = conn.QueryRow(ctx,
		"select content from lyrics where music_id = $1",
		musicID,
	).Scan(&content)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return model.Respond{}, nil
		}
		return model.Respond{}, err
	}

	return content, nil
}
