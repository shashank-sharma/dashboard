package app

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/joho/godotenv"
	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/apis"
	"github.com/pocketbase/pocketbase/cmd"
	"github.com/pocketbase/pocketbase/core"
	"github.com/shashank-sharma/backend/internal/cronjobs"
	"github.com/shashank-sharma/backend/internal/logger"
	"github.com/shashank-sharma/backend/internal/routes"
	"github.com/shashank-sharma/backend/internal/store"
)

type Application struct {
	Server *http.Server
	Pb     *pocketbase.PocketBase
}

func defaultPublicDir() string {
	if strings.HasPrefix(os.Args[0], os.TempDir()) {
		return "./pb_public"
	}

	return filepath.Join(os.Args[0], "../pb_public")
}

func New() *Application {
	pb := pocketbase.New()
	var publicDirFlag string

	pb.RootCmd.PersistentFlags().StringVar(
		&publicDirFlag,
		"publicDir",
		defaultPublicDir(),
		"the directory to serve static files",
	)

	err := godotenv.Load()
	if err != nil {
		fmt.Println("Failed to load environment variables")
	}

	app := &Application{
		Pb: pb,
	}

	pb.OnBeforeServe().Add(func(e *core.ServeEvent) error {
		// Logger need to be initialized inside pocketbase
		// before serve
		logger.RegisterApp(pb)
		dao := pb.Dao()
		store.InitDao(dao)
		app.InitCronjobs()
		app.configureRoutes(e)

		return nil
	})

	return app
}

func (app *Application) configureRoutes(e *core.ServeEvent) {
	e.Router.GET("/api/token", routes.AuthGenerateDevToken, apis.RequireRecordAuth(), apis.ActivityLogger(app.Pb))
	e.Router.POST("/api/track/create", routes.TrackCreateAppItems, apis.RequireRecordAuth(), apis.ActivityLogger(app.Pb))
	e.Router.POST("/api/track", routes.TrackDeviceStatus)
	e.Router.GET("/api/testing", routes.TestHandler)
	e.Router.GET("/stream_mp3", routes.AudioStreamMP3)
	e.Router.POST("/sync/track-items", routes.TrackAppItems, apis.RequireRecordAuth(), apis.ActivityLogger(app.Pb))
	e.Router.GET("/auth/calendar/redirect", routes.CalendarAuthHandler, apis.RequireRecordAuth(), apis.ActivityLogger(app.Pb))
	e.Router.POST("/auth/calendar/callback", routes.CalendarAuthCallback, apis.RequireRecordAuth(), apis.ActivityLogger(app.Pb))
	e.Router.POST("/api/calendar/sync", routes.CalendarSyncHandler, apis.RequireRecordAuth(), apis.ActivityLogger(app.Pb))
}

func (app *Application) InitCronjobs() error {
	cronJobs := []cronjobs.CronJob{
		{
			Name:     "track-device",
			Interval: "*/6 * * * *",
			JobFunc: func() {
				cronjobs.TrackDevices(app.Pb)
			},
			IsActive: true,
		},
	}

	cronjobs.Run(cronJobs)
	return nil
}

func (app *Application) Start() error {

	// register system commands
	app.Pb.RootCmd.AddCommand(cmd.NewAdminCommand(app.Pb))
	app.Pb.RootCmd.AddCommand(cmd.NewServeCommand(app.Pb, true))

	return app.Pb.Execute()
}
