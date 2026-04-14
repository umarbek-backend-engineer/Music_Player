package repository

import (
	"context"
	"database/sql"
	"lyrics-service/internal/repository/posgres"
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
