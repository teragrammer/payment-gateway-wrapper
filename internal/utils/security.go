package utils

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"math/big"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func GeneratePasswordHash(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

func ValidatePassword(hashedPassword, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}

func GenerateRandomString(length int) (string, error) {
	const charset = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"

	// Create a slice to hold the random characters
	var result []byte

	// Generate random characters
	for i := 0; i < length; i++ {
		// Generate a random index in the charset
		idx, err := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		if err != nil {
			return "", err // Return an error if random number generation fails
		}

		// Append the character corresponding to the random index
		result = append(result, charset[idx.Int64()])
	}

	// Return the generated string as a string
	return string(result), nil
}

func GenerateUUID() string {
	return uuid.New().String()
}

func GenerateKey(length int) (string, error) {
	random, err := GenerateRandomString(length)
	if err != nil {
		return "", err
	}

	return GenerateUUID() + "-" + random, nil
}

func SHA256(data string) string {
	hash := sha256.New()
	hash.Write([]byte(data))
	return hex.EncodeToString(hash.Sum(nil))
}
