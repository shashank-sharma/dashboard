package routes

import (
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v5"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/tools/types"
	"github.com/shashank-sharma/backend/internal/models"
	"github.com/shashank-sharma/backend/internal/query"
	"github.com/shashank-sharma/backend/internal/services/fold"
	"github.com/shashank-sharma/backend/internal/util"
)

type FoldOtpAPI struct {
	PhoneNumber string `json:"phone_number"`
	Provider    string `json:"provider"`
}

type FoldVerifyOtpAPI struct {
	PhoneNumber string `json:"phone_number"`
	Otp         string `json:"otp"`
}

func RegisterFoldRoutes(e *core.ServeEvent, foldService *fold.FoldService) {
	e.Router.POST("/api/fold/getotp", func(e *core.RequestEvent) error {
		return FoldGetOtpHandler(foldService, e)
	})
	e.Router.POST("/api/fold/verifyotp", func(e *core.RequestEvent) error {
		return FoldVerifyOtpHandler(foldService, e)
	})
	e.Router.GET("/api/fold/refresh", func(e *core.RequestEvent) error {
		return FoldRefreshTokenHandler(foldService, e)
	})
	e.Router.GET("/api/fold/user", func(e *core.RequestEvent) error {
		return FoldUserHandler(foldService, e)
	})
}

func FoldGetOtpHandler(fs *fold.FoldService, e *core.RequestEvent) error {
	pbToken := e.Request.Header.Get("Authorization")
	userId, _ := util.GetUserId(pbToken)

	foldOtpData := &FoldOtpAPI{}
	if err := e.BindBody(foldOtpData); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid input")
	}

	foldToken := &models.FoldToken{
		User:           userId,
		Phone:          foldOtpData.PhoneNumber,
		UserAgent:      "ua/fold",
		DeviceHash:     uuid.NewString(),
		DeviceLocation: "India",
		DeviceName:     "fold",
		DeviceType:     "Android",
	}
	err := fs.GetOTP(foldToken)
	if err != nil {
		return e.JSON(http.StatusInternalServerError, map[string]string{"error": fmt.Sprintf("error occured: %v", err)})
	}

	if err := query.SaveRecord(foldToken); err != nil {
		return e.JSON(http.StatusInternalServerError, map[string]string{"error": fmt.Sprintf("error saving record: %v", err)})

	}
	return e.JSON(http.StatusOK, map[string]interface{}{"message": "OTP Generated"})
}

func FoldVerifyOtpHandler(fs *fold.FoldService, e *core.RequestEvent) error {
	pbToken := e.Request.Header.Get("Authorization")
	userId, _ := util.GetUserId(pbToken)

	foldData := &FoldVerifyOtpAPI{}
	if err := e.BindBody(foldData); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid input")
	}

	foldToken, err := query.FindByFilter[*models.FoldToken](map[string]interface{}{
		"user": userId,
	})

	if err != nil {
		return e.JSON(http.StatusInternalServerError, map[string]string{"error": "Nothing to verify"})
	}
	verifyResponse, err := fs.VerifyOtp(foldData.Otp, foldToken)
	if err != nil {
		return e.JSON(http.StatusInternalServerError, map[string]string{"error": fmt.Sprintf("error occured: %v", err)})
	}

	foldToken.AccessToken = verifyResponse.Data.AccessToken
	foldToken.RefreshToken = verifyResponse.Data.RefreshToken
	foldToken.Uuid = verifyResponse.Data.UserMeta.UUID
	foldToken.FirstName = verifyResponse.Data.UserMeta.FirstName
	foldToken.LastName = verifyResponse.Data.UserMeta.LastName
	foldToken.Email = verifyResponse.Data.UserMeta.Email

	expiry := types.DateTime{}
	expiry.Scan(verifyResponse.Data.ExpiresAt)
	foldToken.ExpiresAt = expiry

	if err := query.SaveRecord(foldToken); err != nil {
		return e.JSON(http.StatusInternalServerError, map[string]string{"error": "Error updating accessToken and refreshToken"})
	}

	return e.JSON(http.StatusOK, map[string]interface{}{"message": "OTP Verified successfully"})
}

func FoldRefreshTokenHandler(fs *fold.FoldService, e *core.RequestEvent) error {
	pbToken := e.Request.Header.Get("Authorization")
	userId, _ := util.GetUserId(pbToken)

	foldToken, err := query.FindByFilter[*models.FoldToken](map[string]interface{}{
		"user": userId,
	})

	if err != nil {
		return e.JSON(http.StatusInternalServerError, map[string]string{"error": "Nothing to refresh"})
	}

	fs.Refresh(foldToken)
	return e.JSON(http.StatusOK, map[string]interface{}{"message": "Refresh done"})
}

func FoldUserHandler(fs *fold.FoldService, e *core.RequestEvent) error {
	pbToken := e.Request.Header.Get("Authorization")
	userId, _ := util.GetUserId(pbToken)

	foldToken, err := query.FindByFilter[*models.FoldToken](map[string]interface{}{
		"user": userId,
	})

	if err != nil {
		return e.JSON(http.StatusInternalServerError, map[string]string{"error": "Nothing to refresh"})
	}

	fs.Refresh(foldToken)
	fs.User(foldToken)
	return e.JSON(http.StatusOK, map[string]interface{}{"message": "Refresh done"})
}
