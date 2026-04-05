package utils

import (
	"errors"
	"log"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	ErrInvalidFile = errors.New("Invalid file")
	ErrFileTooLarg = errors.New("File too big")
	ErrNotFound    = errors.New("File not found")
)

var mapError = map[error]codes.Code{
	ErrInvalidFile: codes.InvalidArgument,
	ErrFileTooLarg: codes.InvalidArgument,
	ErrNotFound:    codes.NotFound,
}

func MapErrors(err error) error {
	for e, code := range mapError {
		if errors.Is(err, e) {
			return status.Error(code, e.Error())
		}
	}

	log.Println("Internal error: ", err)
	return status.Error(codes.Internal, "Internal Server Error")
}
