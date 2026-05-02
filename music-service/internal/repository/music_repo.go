package repository

import (
	"context"
	"database/sql"
	"fmt"

	pb "github.com/umarbek-backend-engineer/Music_Player/music-service/github.com/umarbek-backend-engineer/Music_Player/music-service/proto/gen"
	"github.com/umarbek-backend-engineer/Music_Player/music-service/internal/model"
	"github.com/umarbek-backend-engineer/Music_Player/music-service/internal/repository/db_connect"
)

func UploadMusicDBHandler(ctx context.Context, user_id, title, filepath string) error {

	conn, err := db_connect.Connect()
	if err != nil {
		return err
	}
	defer conn.Close(ctx)

	msgTag, err := conn.Exec(ctx, "insert into music (user_id, title, filepath) values ($1, $2, $3) on conflict (filepath) do nothing", user_id, title, filepath)
	if err != nil {
		return err
	}

	if msgTag.RowsAffected() == 0 {
		return fmt.Errorf("no rows inserted (conflict or failure)")
	}

	return nil
}

func ListMusicDB(ctx context.Context, user_id string) ([]*pb.MusicItem, error) {
	// connect to database
	conn, err := db_connect.Connect()
	if err != nil {
		return nil, err
	}
	defer conn.Close(ctx)

	// get the music where id matches
	rows, err := conn.Query(ctx, "select id, filename from music where user_id = $1", user_id)
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

func GetMusicIndoFromDB_on_ID(ctx context.Context, user_id, music_id string) (model.Music, error) {
	var music model.Music
	conn, err := db_connect.Connect()
	if err != nil {
		return model.Music{}, err
	}
	defer conn.Close(ctx)

	err = conn.QueryRow(ctx, "select filename, filepath from music where id = $1 and user_id = $2", music_id, user_id).Scan(&music.FileName, &music.FilePath)
	if err != nil {
		return model.Music{}, err
	}
	return music, nil
}
