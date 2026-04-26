package utils

import (
	"errors"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// errors
var (
	ErrExistingEmail    = errors.New("User already exists with the same email")
	ErrInvalidPassword  = errors.New("Invalid Password")
	ErrMetaData         = errors.New("Invalid MetaData")
	ErrUserDoesNotExist = errors.New("User does not exists")
	ErrInvalidToken     = errors.New("Invalid Token")
)

// a map error which will store the variable as a key and code as a value in the map
var mapError = map[error]codes.Code{
	ErrExistingEmail:    codes.InvalidArgument,
	ErrInvalidPassword:  codes.InvalidArgument,
	ErrMetaData:         codes.InvalidArgument,
	ErrUserDoesNotExist: codes.InvalidArgument,
	ErrInvalidToken:     codes.Unauthenticated,
}

func MapErrors(err error) error {
	for e, code := range mapError {
		if errors.Is(err, e) {
			return status.Error(code, e.Error())
		}
	}
	// if the condition of the map error do not match, it will return default which  is internal + error.Error()
	return status.Error(codes.Internal, "Internal Server Error: "+err.Error())
}
