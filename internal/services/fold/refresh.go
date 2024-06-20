package fold

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/pocketbase/pocketbase/tools/types"
	"github.com/shashank-sharma/backend/internal/logger"
	"github.com/shashank-sharma/backend/internal/models"
	"github.com/shashank-sharma/backend/internal/query"
)

type RefreshPayload struct {
	Error
	Refresh string `json:"refresh_token,omitempty"`
}

type RefreshResponse struct {
	Error
	Meta struct {
		RequestID string    `json:"request_id"`
		Timestamp time.Time `json:"timestamp"`
		URI       string    `json:"uri"`
	} `json:"meta"`
	Data struct {
		TokenType    string    `json:"token_type"`
		AccessToken  string    `json:"access_token"`
		RefreshToken string    `json:"refresh_token"`
		ExpiresAt    time.Time `json:"expires_at"`
	} `json:"data"`
}

func (fs *FoldService) RefreshToken(foldToken *models.FoldToken) (RefreshResponse, error) {

	refreshPayloadBuf := new(bytes.Buffer)
	json.NewEncoder(refreshPayloadBuf).Encode(&RefreshPayload{Refresh: foldToken.RefreshToken})

	req, _ := fs.client.NewRequest("POST", "/v1/auth/tokens/refresh", refreshPayloadBuf, foldToken)
	resp, err := fs.client.Do(req)

	if err != nil {
		return RefreshResponse{}, err
	} else {

		logger.LogInfo(fmt.Sprintf("Refresh response status: %+v", resp.StatusCode))

		if resp.StatusCode/100 != 2 {
			return RefreshResponse{}, errors.New(resp.Status)
		}

		data := RefreshResponse{}
		json.NewDecoder(resp.Body).Decode(&data)
		logger.LogInfo(fmt.Sprintf("Verify OTP response: %+v", data))
		return data, nil
	}
}

func (fs *FoldService) Refresh(foldToken *models.FoldToken) {

	if time.Now().After(foldToken.ExpiresAt.Time()) || time.Now().After(foldToken.ExpiresAt.Time()) {
		tokenResponse, err := fs.RefreshToken(foldToken)
		if err != nil {
			logger.LogError("Failed to refresh token for:", foldToken.User)
		}

		foldToken.AccessToken = tokenResponse.Data.AccessToken
		foldToken.RefreshToken = tokenResponse.Data.RefreshToken

		expiresAt := types.DateTime{}
		expiresAt.Scan(tokenResponse.Data.ExpiresAt)
		foldToken.ExpiresAt = expiresAt

		query.SaveRecord(foldToken)
		logger.Debug.Println("Token refreshed")
	} else {
		logger.Debug.Println("Refresh token not needed")
	}
}
