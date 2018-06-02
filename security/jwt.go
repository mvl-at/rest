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
	"strings"
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

func Login(data *JWTData) (bool, string) {
	member := findMember(data.Username, data.Password)

	if member == nil {
		return false, ""
	}

	token := generateToken(member)
	return true, token
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
func generateToken(member *model.Member) string {
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
	return fullToken
}

func fetchMember(id int64) *model.Member {
	member := &model.Member{Id: id}
	database.GenericSingleFetch(member)
	return member
}

func Check(token string) (valid bool, member *model.Member) {
	tokenParts := strings.Split(token, ".")
	payloadJsonByte, _ := base64.URLEncoding.DecodeString(tokenParts[1])
	payload := &JWTPayload{}
	json.Unmarshal(payloadJsonByte, payload)
	member = fetchMember(payload.MemberId)
	valid = true
	if token != generateToken(member) || time.Now().After(payload.Expiration) {
		member = nil
		valid = false
	}
	return
}
