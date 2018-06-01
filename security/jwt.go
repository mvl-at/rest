package security

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"rest/context"
	"rest/database"
	"rest/model"
	"time"
)

type JWTData struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type JWTHeader struct {
	Algorithm string `json:"alg"`
	Type      string `json:"typ"`
}

type JWTPayload struct {
	MemberId   int64     `json:"memberId"`
	Expiration time.Time `json:"exp"`
}

var sessions = map[string]*JWTPayload{}

func Login(data *JWTData) (bool, string) {
	member := findMember(data.Username, data.Password)

	if member == nil {
		return false, ""
	}

	token, payload := generateToken(member)
	sessions[token] = payload
	return true, token
}

func SessionClearer() {
	go clearSessions()
}

func clearSessions() {
	for {
		time.Sleep(time.Minute * 1)
		for k, v := range sessions {

			if time.Now().After(v.Expiration) {
				delete(sessions, k)
			}
		}
	}
}

func findMember(username string, password string) *model.Member {
	members := make([]*model.Member, 0)
	database.GenericFetch(&members)
	for _, v := range members {
		if v.Username == username && v.Password == password {
			return v
		}
	}
	return nil
}

//TODO generate expiration and algorithm type
func generateToken(member *model.Member) (string, *JWTPayload) {
	header := JWTHeader{Algorithm: ""}
	payload := JWTPayload{MemberId: member.Id}
	head, _ := json.Marshal(header)
	pay, _ := json.Marshal(payload)
	encodedHead := base64.URLEncoding.EncodeToString([]byte(head))
	encodedPayload := base64.URLEncoding.EncodeToString([]byte(pay))
	rawToken := encodedHead + "." + encodedPayload
	sig := hmac.New(sha256.New, []byte(context.Conf.JwtSecret))
	sig.Write([]byte(rawToken))
	fullToken := rawToken + "." + hex.EncodeToString(sig.Sum(nil))
	return fullToken, &payload
}
