package security

import (
	"golang.org/x/crypto/bcrypt"
	"rest/database"
	"rest/model"
)

type Credentials struct {
	MemberId int64  `json:"memberId"`
	Username string `json:"username"`
	Password string `json:"password"`
}

func PasswordHash(plainPassword string) string {
	sig, _ := bcrypt.GenerateFromPassword([]byte(plainPassword), bcrypt.DefaultCost)
	return string(sig)
}

func passwordCorrect(plainPassword string, hashedPassword string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(plainPassword)) == nil
}

func UpdateCredentials(credentials *Credentials) {
	member := &model.Member{Id: credentials.MemberId}
	database.Find(member)
	if member != nil {
		member.Username = credentials.Username
		member.Password = PasswordHash(credentials.Password)
		database.Save(member)
	}
}
