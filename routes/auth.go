package routes

import (
	"net/http"

	"github.com/labstack/echo/v5"
	"github.com/pocketbase/dbx"
	"github.com/shashank-sharma/backend/config"
	"github.com/shashank-sharma/backend/logger"
	"github.com/shashank-sharma/backend/util"
)

func AuthGenerateDevToken(c echo.Context) error {
	app := config.GetApp()
	token := c.Request().Header.Get(echo.HeaderAuthorization)
	logger.Debug.Println("token =", token)
	userId, err := util.GetUserId(token)
	if err != nil {
		return c.JSON(http.StatusForbidden, map[string]interface{}{"message": "Failed to fetch id, token misconfigured"})
	}
	record, err := app.Dao().FindFirstRecordByFilter(
		"dev_tokens", "user = {:user}",
		dbx.Params{"user": userId})

	// TODO: Create token if not found and return it
	if err != nil {
		return c.JSON(http.StatusForbidden, map[string]interface{}{"message": "No Dev token found"})
	}

	isActive := false
	devToken := record.Get("token")
	isActive = record.GetBool("is_active")
	if isActive {
		obj := map[string]interface{}{"token": devToken}
		logger.Debug.Println("Successful token sent", obj)
		return c.JSON(http.StatusOK, obj)
	} else {
		return c.JSON(http.StatusForbidden, map[string]interface{}{"message": "Dev token is disabled"})
	}

}
