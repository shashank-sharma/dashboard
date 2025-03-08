package routes

import (
	"net/http"

	"github.com/pocketbase/pocketbase/core"
	"github.com/shashank-sharma/backend/internal/logger"
)

func TestHandler(e *core.RequestEvent) error {
	logger.Debug.Println("Started track")

	return e.JSON(http.StatusOK, map[string]interface{}{"message": "success"})
}
