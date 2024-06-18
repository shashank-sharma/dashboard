package routes

import (
	"net/http"

	"github.com/labstack/echo/v5"
	"github.com/shashank-sharma/backend/internal/logger"
	"github.com/shashank-sharma/backend/internal/models"
	"github.com/shashank-sharma/backend/internal/query"
)

func TestHandler(c echo.Context) error {
	logger.Debug.Println("Started track")

	user, err := query.FindByFilter[*models.Users](map[string]interface{}{
		"id": "",
	})
	if err != nil {
		logger.Debug.Println("err: ", err)
	}
	logger.Debug.Println("user: ", user)
	logger.Debug.Println("Completed")

	logger.Debug.Println("User found:", user)

	if user == nil {
		return c.JSON(http.StatusNotFound, "")
	}

	return c.JSON(http.StatusOK, map[string]interface{}{"message": user.Email})
}
