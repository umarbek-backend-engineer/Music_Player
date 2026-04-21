package utils

import (
	"errors"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	ErrExistingEmail = errors.New("User already exists with the same email")
)

var mapError = map[error]codes.Code{
	ErrExistingEmail: codes.InvalidArgument,
}

func MapErrors(err error) error {
	for e, code := range mapError {
		if errors.Is(err, e) {
			return status.Error(code, e.Error())
		}
	}

	return status.Error(codes.Internal ,"Internal Server Error: " + err.Error())
}
