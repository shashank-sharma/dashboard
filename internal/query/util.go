package query

import (
	"strings"

	"github.com/pocketbase/dbx"
	"github.com/shashank-sharma/backend/internal/models"
)

func structToHashExp(filterStruct map[string]interface{}) dbx.HashExp {
	hashExp := dbx.HashExp{}
	for key, value := range filterStruct {
		hashExp[key] = value
	}
	return hashExp
}

func ValidateDevToken(tokenString string) (string, error) {
	parts := strings.Split(tokenString, ".")
	_, err := FindByFilter[*models.DevToken](map[string]interface{}{
		"user":  parts[0],
		"token": parts[1],
	})
	if err != nil {
		return "", err
	}
	return parts[0], nil
}
