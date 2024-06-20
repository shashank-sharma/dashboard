package fold

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/shashank-sharma/backend/internal/logger"
	"github.com/shashank-sharma/backend/internal/models"
)

type Error struct {
	Error string `json:"error,omitempty"`
}

type LoginPayload struct {
	Error
	Phone   string `json:"phone,omitempty"`
	Channel string `json:"channel,omitempty"`
}

func (fs *FoldService) GetOTP(tokenModel *models.FoldToken) error {
	loginPayloadBuf := new(bytes.Buffer)
	json.NewEncoder(loginPayloadBuf).Encode(&LoginPayload{Phone: tokenModel.Phone, Channel: "sms"})

	req, _ := fs.client.NewRequest("POST", "/v1/auth/otp", loginPayloadBuf, tokenModel)
	resp, err := fs.client.Do(req)

	if err != nil {
		return err
	} else {
		logger.LogInfo(fmt.Sprintf("Login response status: %+v", resp.StatusCode))

		if resp.StatusCode/100 != 2 {
			return errors.New(resp.Status)
		}

		buf := new(bytes.Buffer)
		buf.ReadFrom(resp.Body)
		logger.LogInfo("Login response: " + buf.String())

		return nil
	}
}

type VerifyOtpPayload struct {
	Error
	Phone string `json:"phone,omitempty"`
	Otp   string `json:"otp,omitempty"`
}

type VerifyOtpResponse struct {
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
		NewUser      bool      `json:"new_user"`
		UserID       string    `json:"user_id"`
		UserMeta     struct {
			UUID          string      `json:"uuid"`
			FirstName     string      `json:"first_name"`
			MiddleName    interface{} `json:"middle_name"`
			LastName      string      `json:"last_name"`
			Email         string      `json:"email"`
			EmailVerified bool        `json:"email_verified"`
			GoogleLinked  bool        `json:"google_linked"`
			AppleLinked   bool        `json:"apple_linked"`
			Role          string      `json:"role"`
		} `json:"user_meta"`
	} `json:"data"`
}

func (fs *FoldService) VerifyOtp(otp string, foldToken *models.FoldToken) (VerifyOtpResponse, error) {
	otpVerifyPayloadBuf := new(bytes.Buffer)
	json.NewEncoder(otpVerifyPayloadBuf).Encode(&VerifyOtpPayload{Phone: foldToken.Phone, Otp: otp})

	req, _ := fs.client.NewRequest("POST", "/v1/auth/otp/verify", otpVerifyPayloadBuf, foldToken)
	resp, err := fs.client.Do(req)

	if err != nil {
		return VerifyOtpResponse{}, err
	} else {

		logger.LogInfo(fmt.Sprintf("Otp verify response status: %+v", resp.StatusCode))

		if resp.StatusCode/100 != 2 {
			return VerifyOtpResponse{}, errors.New(resp.Status)
		}

		data := VerifyOtpResponse{}
		json.NewDecoder(resp.Body).Decode(&data)
		logger.LogInfo(fmt.Sprintf("Verify OTP response: %+v", data))
		return data, nil
	}
}
