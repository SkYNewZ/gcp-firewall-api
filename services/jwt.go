package services

import (
	"encoding/base64"
	"encoding/json"
	"strings"

	"github.com/adeo/iwc-gcp-firewall-api/models"
	"github.com/sirupsen/logrus"
)

// JWT describe a Google JSON Web Token
type JWT struct {
	Iss           string `json:"iss"`
	Email         string `json:"email"`
	EmailVerified bool   `json:"email_verified"`
}

// GetUserEmailFromJWT parse the given JWT and return the user email
func GetUserEmailFromJWT(token string) (string, error) {
	logrus.Debugln("Decoding token")

	// Ensure token contains 3 parts
	tokenParts := strings.Split(token, ".")
	if len(tokenParts) != 3 {
		return "", models.NewBadTokenError("Malformed JWT")
	}

	tokenPart := padBase64Input(tokenParts[1])

	data, err := base64.StdEncoding.DecodeString(tokenPart)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"go-err":     err,
			"token-part": tokenParts[1],
		}).Error("Error while base64 decoding token")
		return "", models.NewBadTokenError()
	}

	var t JWT
	err = json.Unmarshal(data, &t)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"go-err":     err,
			"token-data": string(data),
		}).Error("Error while JSON decoding token")
		return "", models.NewBadTokenError()
	}

	if t.Iss != "https://accounts.google.com" {
		logrus.WithFields(logrus.Fields{
			"issuer": t.Iss,
		}).Warningln("Invalid issuer")
		return "", models.NewBadTokenError("Invalid issuer")
	}

	if !t.EmailVerified {
		return "", models.NewBadTokenError("Email not veried")
	}

	logrus.Debugf("User found: %s", t.Email)
	return t.Email, nil
}

// Ensure the given token is base64 decode capable
func padBase64Input(i string) string {
	l := len(i)
	if l%4 == 0 {
		return i
	}

	return i + strings.Repeat("=", 4-(l%4))
}
