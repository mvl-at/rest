package database

import (
	"github.com/mvl-at/model"
	"golang.org/x/crypto/bcrypt"
)

//Defines the credentials with plain password.
type Credentials struct {
	MemberId int64  `json:"memberId"`
	Username string `json:"username"`
	Password string `json:"password"`
}

//Hashes a plaintext password and returns the hash.
func PasswordHash(plainPassword string) string {
	sig, _ := bcrypt.GenerateFromPassword([]byte(plainPassword), bcrypt.DefaultCost)
	return string(sig)
}

//Checks if the given plaintext password matches the hash.
//Returns true if match.
func passwordCorrect(plainPassword string, hashedPassword string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(plainPassword)) == nil
}

//Persists the given credentials to the equivalent member into the 
func UpdateCredentials(credentials *Credentials) {
	member := &model.Member{Id: credentials.MemberId}
	Find(member)
	if member != nil {
		member.Username = credentials.Username
		member.Password = PasswordHash(credentials.Password)
		Save(member)
	}
}
