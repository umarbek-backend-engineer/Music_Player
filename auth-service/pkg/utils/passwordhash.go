package utils

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"

	"golang.org/x/crypto/argon2"
)

func PasswordHash(password string) (string, error) {
	// check if the password is empty or not
	if password == "" {
		return "", errors.New("Invalid Password")
	}
	// create byte slice of 16 byte as a salf
	salt := make([]byte, 16)

	// fill the salt with random numbers
	_, err := rand.Read(salt)
	if err != nil {
		return "", err
	}

	// encode salt+hash using argon2 (cannot be decoded) for more security
	hash := argon2.IDKey([]byte(password), salt, 1, 64*1024, 4, 32)

	// these salt and hash based64 will be saved in database
	saltBase64 := base64.StdEncoding.EncodeToString(salt)
	hashBase64 := base64.StdEncoding.EncodeToString(hash)

	// make them together, in db they will be saved as salt.hash and when checking the password the system will extract salt and decoded
	// with decoded salt it will created the exact same hash using argon2 algorith and copare the hash with the dbhash
	encodedHash := fmt.Sprintf("%s.%s", saltBase64, hashBase64)
	return encodedHash, nil
}
