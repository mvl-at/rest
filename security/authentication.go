package security

import (
	"encoding/hex"
	"golang.org/x/crypto/bcrypt"
)

func passwordHash(plainPassword string) string {
	sig, _ := bcrypt.GenerateFromPassword([]byte(plainPassword), bcrypt.DefaultCost)
	return hex.EncodeToString(sig)
}

func passwordCorrect(plainPassword string, hashedPassword string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(plainPassword)) == nil
}
