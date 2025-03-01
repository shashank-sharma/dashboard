package routes

import (
	"net/http"

	"github.com/pocketbase/pocketbase/core"
	"github.com/shashank-sharma/backend/internal/models"
	"github.com/shashank-sharma/backend/internal/query"
	"github.com/shashank-sharma/backend/internal/util"
)

func RegisterCredentialRoutes(e *core.ServeEvent) {
	e.Router.GET("/api/token", AuthGenerateDevToken)
}

func AuthGenerateDevToken(e *core.RequestEvent) error {
	token := e.Request.Header.Get("Authorization")
	userId, err := util.GetUserId(token)
	if err != nil {
		return e.JSON(http.StatusForbidden, map[string]interface{}{"message": "Failed to fetch id, token misconfigured"})
	}

	record, err := query.FindByFilter[*models.DevToken](map[string]interface{}{
		"user": userId,
	})

	// TODO: Create token if not found and return it
	if err != nil {
		return e.JSON(http.StatusNotFound, map[string]interface{}{"message": "No Dev token found"})
	}

	if record.IsActive {
		obj := map[string]interface{}{"token": record.Token}
		return e.JSON(http.StatusOK, obj)
	} else {
		return e.JSON(http.StatusForbidden, map[string]interface{}{"message": "Dev token is disabled"})
	}

}
