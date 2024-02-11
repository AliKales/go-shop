package tokengenerator

import (
	"crypto/rand"
	"encoding/hex"
)

func GenerateSecureToken(length int) string {
	b := make([]byte, length)
	if _, err := rand.Read(b); err != nil {
		return ""
	}
	return hex.EncodeToString(b)
}

func GenerateUserToken() string {
	return GenerateSecureToken(50)
}

func GenerateUserRefreshToken() string {
	return GenerateSecureToken(100)
}