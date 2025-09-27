package common

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

var JwtSecret = []byte("access-control-management-secret-key")

const workFactor = 12 // ค่ามาตรฐานทั่วไปคือ 10-12, ค่าที่สูงขึ้นคือปลอดภัยขึ้น

func HashPassword(password string) (string, error) {
	pwdBytes := []byte(password)
	hashedBytes, err := bcrypt.GenerateFromPassword(pwdBytes, workFactor)
	if err != nil {
		return "", fmt.Errorf("failed to hash password: %w", err)
	}

	return string(hashedBytes), nil
}

// --- Verification Function ---

// CompareHashAndPassword checks if a plain-text password matches the stored bcrypt hash.
func VerifyHashPassword(hashedPassword string, password string) bool {
	err := bcrypt.CompareHashAndPassword(
		[]byte(hashedPassword),
		[]byte(password),
	)
	return err == nil
}
