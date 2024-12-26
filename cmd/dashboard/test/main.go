package main

import (
	"github.com/shashank-sharma/backend/internal/config"
	"github.com/shashank-sharma/backend/internal/logger"
	"github.com/shashank-sharma/backend/internal/routes"
	"github.com/shashank-sharma/backend/internal/store"

	"github.com/pocketbase/pocketbase/core"
)

func main() {

	config.Init()
	app := config.GetApp()
	app.Logger().Error("Debug message")

	app.OnServe().BindFunc(func(e *core.ServeEvent) error {
		logger.RegisterApp(app)
		store.InitApp(app)

		logger.LogError(
			"Debug message with attributes!",
			"name", "Johnnyyyy1 Doe",
			"id", 123,
		)

		e.Router.GET("/api/testing",
			routes.TestHandler)
		return nil
	})
}
