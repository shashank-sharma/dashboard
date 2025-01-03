package util

import (
	"encoding/base64"
	"encoding/json"
	"strings"

	_ "github.com/mattn/go-sqlite3"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/tools/security"
	"github.com/shashank-sharma/backend/internal/logger"
)

func GenerateRandomId() string {
	return security.RandomStringWithAlphabet(core.DefaultIdLength, core.DefaultIdAlphabet)
}

func GetUserId(tokenString string) (string, error) {
	// Split the token into header, payload, and signature
	parts := strings.Split(tokenString, ".")

	// Decode the payload (no signature verification)
	payload, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		logger.Error.Println("Error decoding payload:", err)
		return "", err
	}

	// Unmarshal the payload
	var claims map[string]interface{}
	err = json.Unmarshal(payload, &claims)
	if err != nil {
		logger.Error.Println("Error unmarshalling payload:", err)
		return "", err
	}

	// Access the claims
	logger.Debug.Println("User ID:", claims["id"])
	logger.Debug.Println("Username:", claims["username"])
	return claims["id"].(string), nil
}
