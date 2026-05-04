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

	// connect to database
	conn, err := db_connect.Connect()
	if err != nil {
		return err
	}
	defer conn.Close(ctx)

	// store the meta information inside the db
	msgTag, err := conn.Exec(ctx, "insert into music (user_id, title, filepath) values ($1, $2, $3) on conflict (filepath) do nothing", user_id, title, filepath)
	if err != nil {
		return err
	}

	// check it anything was inserted
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
	rows, err := conn.Query(ctx, "select id, title from music where user_id = $1", user_id)
	if err != nil {
		return nil, err
	}

	var musics []*pb.MusicItem

	// going through each row and saing each row inside music var and append to musics slice of pb.MusicItem
	for rows.Next() {
		var music pb.MusicItem
		err = rows.Scan(&music.Id, &music.Title)
		if err != nil && err != sql.ErrNoRows {
			return nil, err
		}

		musics = append(musics, &music)
	}
	return musics, nil
}

func GetMusicIndoFromDB_on_ID(ctx context.Context, music_id string) (model.Music, error) {
	var music model.Music
	conn, err := db_connect.Connect()
	if err != nil {
		return model.Music{}, err
	}
	defer conn.Close(ctx)

	err = conn.QueryRow(ctx, "select title, filepath from music where id = $1", music_id).Scan(&music.FileName, &music.FilePath)
	if err != nil {
		return model.Music{}, err
	}
	return music, nil
}

func GetPublicCrudHandler(ctx context.Context, user_id string) (*pb.PublicMusicResponse, error) {
	conn, err := db_connect.Connect()
	if err != nil {
		return nil, err
	}
	defer conn.Close(ctx)

	// give the command to the database
	rows, err := conn.Query(ctx, "select id, user_id, title, is_public from music where user_id = $1 and is_public != false", user_id)

	// slice of music
	music := &pb.PublicMusicResponse{}
	// loop untill the completion command
	for rows.Next() {
		var single_music pb.PublicMusic
		// assign the information recieved from music to variable single_music
		err = rows.Scan(
			&single_music.MusicId,
			&single_music.UserId,
			&single_music.Title,
			&single_music.IsPublic,
		)

		// add single_music to sliece
		music.Music = append(music.Music, &single_music)
	}
	return music, nil
}

func MakeMusicVisibalCrudHandler(ctx context.Context, is_public bool, user_id, music_id string) error {
	conn, err := db_connect.Connect()
	if err != nil {
		return err
	}
	defer conn.Close(ctx)

	// give the command to database
	_, err = conn.Exec(ctx, "update music set is_public = $1 where user_id = $2 and id = $3", is_public, user_id, music_id)
	if err != nil {
		return err
	}

	return nil
}
