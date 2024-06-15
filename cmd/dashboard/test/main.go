package main

import (
	"log"

	"github.com/shashank-sharma/backend/config"
	"github.com/shashank-sharma/backend/logger"
	"github.com/shashank-sharma/backend/routes"
	"github.com/shashank-sharma/backend/store"

	"github.com/pocketbase/pocketbase/core"
)

func main() {

	config.Init()
	app := config.GetApp()
	app.Logger().Error("Debug message")

	app.OnBeforeServe().Add(func(e *core.ServeEvent) error {
		logger.RegisterApp(app)
		dao := app.Dao()
		// logger.Debug.Println("Dao is: ", dao)
		store.InitDao(dao)

		logger.LogError(
			"Debug message with attributes!",
			"name", "Johnnyyyy1 Doe",
			"id", 123,
		)

		e.Router.GET("/api/testing",
			routes.TestHandler)
		return nil
	})

	if err := app.Start(); err != nil {
		log.Fatal(err)
	}
}
