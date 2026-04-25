package utils

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"errors"
	"fmt"
	"strings"

	"golang.org/x/crypto/argon2"
)

// HashPassword takes a plain-text password, hashes it securely using Argon2,
// and returns a single encoded string containing both the salt and the hash.
// Format of returned string: "saltBase64.hashBase64"
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

	// Combine salt and hash into one string separated by "."
	// Example: "salt123.hash456"
	encodedHash := fmt.Sprintf("%s.%s", saltBase64, hashBase64)
	return encodedHash, nil
}

func VerifyPassword(password, dbpassword string) error {
	// split the salt and the hashpassword
	parts := strings.Split(dbpassword, ".")
	if len(parts) != 2 {
		return fmt.Errorf("Failed to get parts of the dbpassword: len != 2")
	}

	// assign them into variables
	saltbase64 := parts[0]
	hashBase64 := parts[1]

	// decode the salt and the hash to compare them with the salt + "." + password
	decodedsalt, err := base64.StdEncoding.DecodeString(saltbase64)
	if err != nil {
		return err
	}
	decodedhash, err := base64.StdEncoding.DecodeString(hashBase64)
	if err != nil {
		return err
	}

	// encode the request password with the same algorithm argon2 and add the same salt
	passwordHash := argon2.IDKey([]byte(password), []byte(decodedsalt), 1, 64*1024, 4, 32)

	// checking if the passwordhash and dbhash are the same
	// extra security
	if len(passwordHash) != len(decodedhash) {
		return ErrInvalidPassword
	}

	// this function will return 1 if the password do match
	if subtle.ConstantTimeCompare(passwordHash, decodedhash) == 1 {
		return nil
	}

	// return password failed if the password incorrect
	return ErrInvalidPassword
}
