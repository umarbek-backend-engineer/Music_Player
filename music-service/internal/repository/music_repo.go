package repository

import (
	"context"
	"database/sql"
	"music-service/internal/repository/db_connect.go"
	pb "music-service/proto/gen"
)

func UploadMusicDBHandler(ctx context.Context, filename, filepath string) error {

	conn, err := db_connect.Connect()
	if err != nil {
		return err
	}
	defer conn.Close(ctx)

	_, err = conn.Exec(ctx, "insert into musics (filename, filepath) values ($1, $2)", filename, filepath)
	if err != nil {
		return err
	}
	return nil
}

func ListMusicDB(ctx context.Context) ([]*pb.MusicItem, error) {
	conn, err := db_connect.Connect()
	if err != nil {
		return nil, err
	}
	defer conn.Close(ctx)

	rows, err := conn.Query(ctx, "select id, filename from musics")
	if err != nil {
		return nil, err
	}

	var musics []*pb.MusicItem

	// going through each row and saing each row inside music var and append to musics slice of pb.MusicItem
	for rows.Next() {
		var music pb.MusicItem
		err = rows.Scan(&music.Id, &music.Filename)
		if err != nil && err != sql.ErrNoRows {
			return nil, err
		}

		musics = append(musics, &music)
	}
	return musics, nil
}
