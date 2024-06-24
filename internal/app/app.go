package app

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v5"
	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/apis"
	"github.com/pocketbase/pocketbase/cmd"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/plugins/migratecmd"
	"github.com/shashank-sharma/backend/internal/cronjobs"
	"github.com/shashank-sharma/backend/internal/logger"
	"github.com/shashank-sharma/backend/internal/routes"
	"github.com/shashank-sharma/backend/internal/services/calendar"
	"github.com/shashank-sharma/backend/internal/services/fold"
	"github.com/shashank-sharma/backend/internal/store"
)

type Application struct {
	Server          *http.Server
	Pb              *pocketbase.PocketBase
	FoldService     *fold.FoldService
	CalendarService *calendar.CalendarService
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

	// Define all the service here
	foldService := fold.NewFoldService("https://api.fold.money/api")
	calendarService := calendar.NewCalendarService()

	app := &Application{
		Pb:              pb,
		FoldService:     foldService,
		CalendarService: calendarService,
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
	e.Router.GET("/auth/calendar/redirect", func(c echo.Context) error {
		return routes.CalendarAuthHandler(app.CalendarService, c)
	}, apis.RequireRecordAuth(), apis.ActivityLogger(app.Pb))
	e.Router.POST("/auth/calendar/callback", func(c echo.Context) error {
		return routes.CalendarAuthCallback(app.CalendarService, c)
	}, apis.RequireRecordAuth(), apis.ActivityLogger(app.Pb))
	e.Router.POST("/api/calendar/sync", func(c echo.Context) error {
		return routes.CalendarSyncHandler(app.CalendarService, c)
	}, apis.RequireRecordAuth(), apis.ActivityLogger(app.Pb))
	e.Router.POST("/api/fold/getotp", func(c echo.Context) error {
		return routes.FoldGetOtpHandler(app.FoldService, c)
	}, apis.RequireRecordAuth(), apis.ActivityLogger(app.Pb))
	e.Router.POST("/api/fold/verifyotp", func(c echo.Context) error {
		return routes.FoldVerifyOtpHandler(app.FoldService, c)
	}, apis.RequireRecordAuth(), apis.ActivityLogger(app.Pb))
	e.Router.GET("/api/fold/refresh", func(c echo.Context) error {
		return routes.FoldRefreshTokenHandler(app.FoldService, c)
	}, apis.RequireRecordAuth(), apis.ActivityLogger(app.Pb))
	e.Router.GET("/api/fold/user", func(c echo.Context) error {
		return routes.FoldUserHandler(app.FoldService, c)
	}, apis.RequireRecordAuth(), apis.ActivityLogger(app.Pb))
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

	isGoRun := strings.HasPrefix(os.Args[0], os.TempDir())
	logger.Debug.Println("isGoRun:", isGoRun)

	migratecmd.MustRegister(app.Pb, app.Pb.RootCmd, migratecmd.Config{
		// enable auto creation of migration files when making collection changes in the Admin UI
		// (the isGoRun check is to enable it only during development)
		Automigrate: isGoRun,
	})

	return app.Pb.Execute()
}
