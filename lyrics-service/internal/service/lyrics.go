package service

import (
	"bytes"
	"context"
	"io"
	"log"

	pb "github.com/umarbek-backend-engineer/Music_Player/lyrics-service/github.com/umarbek-backend-engineer/Music_Player/lyrics-service/proto/gen"
	"github.com/umarbek-backend-engineer/Music_Player/lyrics-service/internal/repository"
	"github.com/umarbek-backend-engineer/Music_Player/lyrics-service/pkg/utils"
	"google.golang.org/grpc"
)

// This is add lyrics rpc
// first it will ceck if the music exists by check the same names in data base. If yes the function will not create another lyrics, else it will.
// if there is not same music name in data base, it will connect to music-service and pull the stream music rpc
// the rpc will return music bytes and the bytes are send to wisper
// recieved lyrics with timestamp are saved in db

func (s *Server) AddLyrics(ctx context.Context, req *pb.AddLyricsRequest) (*pb.Empty, error) {

	exists, err := repository.Is_music_lyric_exists(ctx, req.Text)
	if err != nil {
		return nil, err
	}
	if exists {
		return &pb.Empty{}, nil
	}

	// connecting  to grpc server to get the music it self
	conn, err := grpc.Dial("music-service:50051", grpc.WithInsecure())
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	client := pb.NewMusicServiceClient(conn)

	// requesting for stream rpc
	stream, err := client.StreamMusic(ctx, &pb.StreamRequest{
		Id: req.MusicId,
	})

	if err != nil {
		return nil, utils.MapError(err)
	}
	var filename string
	var buffer bytes.Buffer
	for {
		//recieving the music chunks
		res, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, utils.MapError(err)
		}

		if res.Name != "" {
			filename = res.Name
		}

		// writing in buffer
		_, err = buffer.Write(res.Content)
		if err != nil {
			log.Println("Error writing the file: ", err)
			break
		}
	}

	LyricsResp, err := utils.SendToWisper(buffer.Bytes(), filename)
	if err != nil {
		return nil, utils.MapError(err)
	}

	err = repository.SaveLyrics(ctx, req.MusicId, req.Text, LyricsResp)
	if err != nil {
		return nil, utils.MapError(err)
	}

	return &pb.Empty{}, nil
}

func (s *Server) GetLyrics(ctx context.Context, req *pb.GetLyricsRequest) (*pb.LyricsResponse, error) {

	// get the lyrics in models struct (json)
	resp, err := repository.GetLyricsByMusicID(ctx, req.MusicId)
	if err != nil {
		return nil, utils.MapError(err)
	}

	// initialize lyricsStruct of protobuffer
	lyrics := make([]*pb.LyricsStruct, 0, len(resp.Lyrics))

	// model struct (json) to pb struct
	for _, seg := range resp.Lyrics {
		lyrics = append(lyrics, &pb.LyricsStruct{
			Start: float32(seg.Start),
			End:   float32(seg.End),
			Text:  seg.Text,
		})
	}

	// response to the GATEWAY
	return &pb.LyricsResponse{
		Language: resp.Language,
		Lyrics:   lyrics,
	}, nil
}
