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

//Defines the JWT data which has to be send to this server.
type JWTData struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

//Defines the JWT header.
type JWTHeader struct {
	Algorithm string `json:"alg"`
	Type      string `json:"typ"`
}

//Defines the JWT payload which is used to communicate between server and client.
type JWTPayload struct {
	MemberId   int64     `json:"memberId"`
	Expiration time.Time `json:"exp"`
}

//Does the login process.
//Returns true, if the login was successful.
//Returns the JWT for future request which require certain roles or permissions.
func Login(data *JWTData) (bool, string) {
	member := findMember(data.Username, data.Password)

	if member == nil {
		return false, ""
	}

	token := generateToken(member)
	return true, token
}

//Finds a member with the given credentials.
//Return nil if nothing was found.
func findMember(username string, password string) *model.Member {
	members := make([]*model.Member, 0)
	database.FindAll(&members)
	for _, v := range members {
		if v.Username == username && passwordCorrect(password, v.Password) && v.LoginAllowed && !v.Deleted {
			return v
		}
	}
	return nil
}

//Returns a JWT for a given Member.
func generateToken(member *model.Member) string {
	header := JWTHeader{Algorithm: ""}
	payload := JWTPayload{MemberId: member.Id, Expiration: time.Now().Add(time.Minute * time.Duration(context.Conf.JwtExpiration))}
	head, _ := json.Marshal(header)
	pay, _ := json.Marshal(payload)
	encodedHead := base64.URLEncoding.EncodeToString([]byte(head))
	encodedPayload := base64.URLEncoding.EncodeToString([]byte(pay))
	rawToken := encodedHead + "." + encodedPayload
	fullToken := rawToken + "." + hash(rawToken)
	return fullToken
}

//Hashes a JWT.
func hash(rawToken string) string {
	sig := hmac.New(sha256.New, []byte(context.Conf.JwtSecret))
	sig.Write([]byte(rawToken))
	return hex.EncodeToString(sig.Sum(nil))
}

//Shortcut to find a member.
func fetchMember(id int64) *model.Member {
	member := &model.Member{Id: id}
	database.Find(member)
	return member
}

//Checks, if the given JWT is valid.
//Returns true if it is valid.
//Returns the member which belongs to the JWT. If the JWT is not valid, the member will be nil.
func Check(token string) (valid bool, member *model.Member) {
	tokenParts := strings.Split(token, ".")
	valid = false

	if len(tokenParts) != 3 || hash(tokenParts[0]+"."+tokenParts[1]) != tokenParts[2] {
		return
	}

	payloadJsonByte, _ := base64.URLEncoding.DecodeString(tokenParts[1])
	payload := &JWTPayload{}
	json.Unmarshal(payloadJsonByte, payload)

	if time.Now().After(payload.Expiration) {
		return
	}

	member = fetchMember(payload.MemberId)

	if member.Deleted || !member.LoginAllowed {
		return
	}

	valid = true
	return
}
