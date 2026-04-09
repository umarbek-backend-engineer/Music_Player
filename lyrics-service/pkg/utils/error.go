package utils

import (
	"log"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func MapError(err error) error {
	log.Println("Internal Error: ", err)
	return status.Error(codes.Internal, "Internal Server Error`")
}
