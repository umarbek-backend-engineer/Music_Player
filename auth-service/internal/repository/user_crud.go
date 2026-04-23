package repository

import (
	"context"

	pb "github.com/umarbek-backend-engineer/Music_Player/github.com/umarbek-backend-engineer/Music_Player/auth-service/proto/gen"
	"github.com/umarbek-backend-engineer/Music_Player/internal/repository/postgres"
)

func RegisterDBCrud(ctx context.Context, req *pb.RegisterRequest) (string, error) {
	// connect to database
	conn, err := postgres.Connect()
	if err != nil {
		return "", err
	}
	// you always have to close the database client(conn)
	defer conn.Close(ctx)

	var id string
	// the query will save the information into the table users in database and return id of the inserted row
	// the id will be stored in id variable
	err = conn.QueryRow(ctx, "insert into users (name, lastname, email, role, password) values ($1,$2,$3,$4,$5, $6) returning id",
		req.Name,
		req.Lastname,
		req.Email,
		req.Role,
		req.Password,
	).Scan(&id)
	if err != nil {
		return "", err
	}
	return id, nil
}
