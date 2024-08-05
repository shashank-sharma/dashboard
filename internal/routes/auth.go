package routes

import (
	"net/http"

	"github.com/labstack/echo/v5"
	"github.com/shashank-sharma/backend/internal/logger"
	"github.com/shashank-sharma/backend/internal/models"
	"github.com/shashank-sharma/backend/internal/query"
	"github.com/shashank-sharma/backend/internal/util"
)

func AuthGenerateDevToken(c echo.Context) error {
	token := c.Request().Header.Get(echo.HeaderAuthorization)
	logger.Debug.Println("token =", token)
	userId, err := util.GetUserId(token)
	if err != nil {
		return c.JSON(http.StatusForbidden, map[string]interface{}{"message": "Failed to fetch id, token misconfigured"})
	}
	logger.Debug.Println("userId =", userId)

	record, err := query.FindByFilter[*models.DevToken](map[string]interface{}{
		"user": userId,
	})

	// TODO: Create token if not found and return it
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]interface{}{"message": "No Dev token found"})
	}

	if record.IsActive {
		obj := map[string]interface{}{"token": record.Token}
		logger.Debug.Println("Successful token sent", obj)
		return c.JSON(http.StatusOK, obj)
	} else {
		return c.JSON(http.StatusForbidden, map[string]interface{}{"message": "Dev token is disabled"})
	}

}
