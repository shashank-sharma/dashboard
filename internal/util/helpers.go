package util

import (
	"encoding/base64"
	"encoding/json"
	"strings"

	_ "github.com/mattn/go-sqlite3"
	pbModels "github.com/pocketbase/pocketbase/models"
	"github.com/pocketbase/pocketbase/tools/security"
	"github.com/shashank-sharma/backend/internal/logger"
	"github.com/shashank-sharma/backend/internal/models"
	"github.com/shashank-sharma/backend/internal/query"
)

func GenerateRandomId() string {
	return security.RandomStringWithAlphabet(pbModels.DefaultIdLength, pbModels.DefaultIdAlphabet)
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

func ValidateDevToken(tokenString string) (string, error) {
	parts := strings.Split(tokenString, ".")
	_, err := query.FindByFilter[*models.DevToken](map[string]interface{}{
		"user":  parts[0],
		"token": parts[1],
	})
	if err != nil {
		return "", err
	}
	return parts[0], nil
}
