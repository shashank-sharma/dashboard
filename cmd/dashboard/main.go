package main

import (
	"log"

	"github.com/shashank-sharma/backend/config"
	"github.com/shashank-sharma/backend/cronjobs"
	"github.com/shashank-sharma/backend/logger"
	"github.com/shashank-sharma/backend/routes"
	"github.com/shashank-sharma/backend/store"

	"github.com/pocketbase/pocketbase/apis"
	"github.com/pocketbase/pocketbase/core"
)

func main() {

	config.Init()
	app := config.GetApp()

	app.OnBeforeServe().Add(func(e *core.ServeEvent) error {
		logger.RegisterApp(app)

		dao := app.Dao()
		store.InitDao(dao)
		cronjobs.InitCronjobs(app)

		// serves static files from the provided public dir (if exists)
		// e.Router.GET("/*", apis.StaticDirectoryHandler(os.DirFS(publicDirFlag), true))

		e.Router.GET("/api/token", routes.AuthGenerateDevToken, apis.RequireRecordAuth(), apis.ActivityLogger(app))

		e.Router.POST("/api/track/create", routes.TrackCreateAppItems, apis.RequireRecordAuth(), apis.ActivityLogger(app))

		e.Router.POST("/api/track", routes.TrackDeviceStatus)

		e.Router.GET("/api/testing",
			routes.TestHandler)

		e.Router.GET("/stream_mp3", routes.AudioStreamMP3)

		e.Router.POST("/sync/track-items",
			routes.TrackAppItems,
			apis.RequireRecordAuth(),
			apis.ActivityLogger(app))

		e.Router.GET("/auth/calendar/redirect", routes.CalendarAuthHandler, apis.RequireRecordAuth(), apis.ActivityLogger(app))
		e.Router.POST("/auth/calendar/callback", routes.CalendarAuthCallback, apis.RequireRecordAuth(), apis.ActivityLogger(app))
		e.Router.POST("/api/calendar/sync", routes.CalendarSyncHandler, apis.RequireRecordAuth(), apis.ActivityLogger(app))

		return nil
	})

	if err := app.Start(); err != nil {
		log.Fatal(err)
	}
}
