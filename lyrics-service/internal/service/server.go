package service

import lyricspb "lyrics-service/proto/gen"

type Server struct {
	lyricspb.UnimplementedLyricsServiceServer
}
