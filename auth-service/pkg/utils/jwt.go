package utils

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/umarbek-backend-engineer/Music_Player/internal/config"
)

// create access JWT token using id
func GenerateAccessJWT(id, role string) (string, error) {
	// loading config of the jwt
	cgf := config.Load()

	// parsing the duration of the access token which is stored in .env file as example(10m)
	// here it will make "10m" = 10 minutes
	duration, err := time.ParseDuration(cgf.ACC_JWT_exp)
	if err != nil {
		return "", err
	}

	// jwt content
	claims := jwt.MapClaims{
		"user_id":   id,
		"role":      role,
		"exp":       jwt.NewNumericDate(time.Now().Add(duration)).Unix(),
		"issued_at": jwt.NewNumericDate(time.Now()).Unix(),
	}

	// making token with all the claims + defining the signing method
	token := jwt.NewWithClaims(jwt.SigningMethodPS256, claims)

	// signing the token with the secret key
	signedToken, err := token.SignedString([]byte(cgf.JWT_key))
	if err != nil {
		return "", err
	}

	// return the complete token
	return signedToken, nil
}

// this function will generate a slice of bytes 256bite, which will be treated as refresh token
func GenerateRefreshTokne() (string, error) {
	// a slice of bytes
	token := make([]byte, 32)

	// fill the slice of bytes using rand.Read with random value and it makes our token
	_, err := rand.Read(token)
	if err != nil {
		return "", err
	}

	// return the token and nill for an error
	return base64.URLEncoding.EncodeToString(token), nil

}

// this function will encode the token
func HashToken(token string) string {

	hashToken := sha256.Sum256([]byte(token))
	return hex.EncodeToString(hashToken[:])
}
