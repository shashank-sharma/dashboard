package middleware

import (
	"github.com/pocketbase/pocketbase/core"
	"github.com/shashank-sharma/backend/internal/logger"
)


func RegisterGlobalMiddleware(e *core.RequestEvent) error {
	logger.Debug.Println("Registering global middleware")
	return e.Next()
}