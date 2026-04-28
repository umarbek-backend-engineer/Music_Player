package repository

import (
	"context"
	"database/sql"
	"time"

	pb "github.com/umarbek-backend-engineer/Music_Player/github.com/umarbek-backend-engineer/Music_Player/auth-service/proto/gen"
	"github.com/umarbek-backend-engineer/Music_Player/internal/repository/postgres"
	"github.com/umarbek-backend-engineer/Music_Player/pkg/utils"
)

// crud operation of the register method it will saved the req information inside the database and will return id
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
	err = conn.QueryRow(ctx, "insert into users (name, lastname, email, role, password) values ($1,$2,$3,$4,$5) returning id",
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

// the crud function of the logout method which will delete the session row where passing refresh token matches
func LogoutCrud(ctx context.Context, hashtoken string) error {
	// create the database client and close it after its usage
	conn, err := postgres.Connect()
	if err != nil {
		return err
	}
	defer conn.Close(ctx)

	var exists bool
	// check if it exists and valid, and return 1 which is true
	err = conn.QueryRow(ctx, "delete from sessions where refresh_token = $1", hashtoken).Scan(&exists)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil
		}
		return err
	}
	return nil
}

func ResetPasswordCrud(ctx context.Context, user_id string, currentPassword string, newpassword string) (string, error) {

	var dbPassword string
	var role string
	// connect Database
	conn, err := postgres.Connect()
	if err != nil {
		return "", err
	}
	// close the client after usage
	defer conn.Close(ctx)

	// Get user password from DB by ID
	err = conn.QueryRow(ctx, "select password from users where id = $1", user_id).Scan(&dbPassword)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", utils.ErrUserDoesNotExist
		}
		return "", err
	}

	// if db and current passwords match it will not return err
	err = utils.VerifyPassword(currentPassword, dbPassword)

	// hash the new password to save into database
	hashedPassword, err := utils.PasswordHash(newpassword)
	if err != nil {
		return "", err
	}

	// update the password and set time now for password_changed_at
	err = conn.QueryRow(ctx, "update users set password = $1, password_changed_at = now() where id = $2 returning role", hashedPassword, user_id).Scan(&role)
	if err != nil {
		return "", err
	}

	// delete all sessions based on user_id
	_, err = conn.Exec(ctx, "delete from sessions where user_id = $1", user_id)
	if err != nil {
		if err == sql.ErrNoRows {
			return role, nil
		}
		return "", err
	}
	return role, nil
}

func RefreshTokenCrud(ctx context.Context, hashed_refresh_token string) (string, string, error) {
	// create database client
	conn, err := postgres.Connect()
	if err != nil {
		return "", "", err
	}
	// close the client after usage
	defer conn.Close(ctx)

	// initialize the needed variable
	var user_id string
	var role string
	var expires_at time.Time
	var revoked bool

	// find the session
	err = conn.QueryRow(ctx, "select s.user_id, s.expires_at, s.revoked, u.role from sessions s join users u on s.user_id = u.id where s.token_hash = $1", hashed_refresh_token).Scan(
		&user_id,
		&expires_at,
		&role,
		&revoked,
	)
	if err != nil {
		return "", "", err
	}

	// Validate session:
	// - if revoked = true → error
	// - if expires_at < now() → error

	if revoked && time.Now().Unix() > expires_at.Unix() {
		return "", "", err
	}
	return user_id, role, nil
}

// delete accoutn crud operations function, it will delete user row where id matrches
func DeleAccountCrud(ctx context.Context, id string) error {
	// create db client
	conn, err := postgres.Connect()
	if err != nil {
		return err
	}
	// close the client after usage
	defer conn.Close(ctx)

	// Query which will delete the row of the user with that id
	_, err = conn.Exec(ctx, "delete from users where id = $1", id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil
		}
		return err
	}
	return nil
}

// crud operations of the logIn method, from the given email it will return user_id, role, and saved password
func LogInCrud(ctx context.Context, email string) (string, string, string, error) {
	conn, err := postgres.Connect()
	if err != nil {
		return "", "", "", err
	}
	// close the client adter usage
	defer conn.Close(ctx)

	// initialize variable to store data from database
	var id string
	var role string
	var dbpassword string

	// query to get id, role and password based on email
	err = conn.QueryRow(ctx, "select id, role, password from users where email = $1", email).Scan(&id, &role, &dbpassword)
	if err != nil {
		return "", "", "", err
	}
	return id, role, dbpassword, nil
}
