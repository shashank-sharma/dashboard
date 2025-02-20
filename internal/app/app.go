package app

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/joho/godotenv"
	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/cmd"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/plugins/migratecmd"
	"github.com/shashank-sharma/backend/internal/config"
	"github.com/shashank-sharma/backend/internal/cronjobs"
	"github.com/shashank-sharma/backend/internal/logger"
	"github.com/shashank-sharma/backend/internal/routes"
	"github.com/shashank-sharma/backend/internal/services/calendar"
	"github.com/shashank-sharma/backend/internal/services/fold"
	"github.com/shashank-sharma/backend/internal/services/mail"
	"github.com/shashank-sharma/backend/internal/store"
)

type Application struct {
	Server          *http.Server
	Pb              *pocketbase.PocketBase
	FoldService     *fold.FoldService
	CalendarService *calendar.CalendarService
	MailService     *mail.MailService
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
	mailService := mail.NewMailService()

	app := &Application{
		Pb:              pb,
		FoldService:     foldService,
		CalendarService: calendarService,
		MailService:     mailService,
	}

	pb.OnServe().BindFunc(func(e *core.ServeEvent) error {
		// Logger need to be initialized inside pocketbase
		// before serve
		logger.RegisterApp(pb)
		store.InitApp(pb)
		config.Init(app.Pb)
		app.InitCronjobs()
		app.configureRoutes(e)

		return e.Next()
	})

	app.registerHooks()

	return app
}

func (app *Application) configureRoutes(e *core.ServeEvent) {
	e.Router.GET("/api/token", routes.AuthGenerateDevToken)
	e.Router.POST("/api/track/create", routes.TrackCreateAppItems)
	e.Router.POST("/api/track", routes.TrackDeviceStatus)
	e.Router.GET("/api/track/getapp", routes.GetCurrentApp)
	e.Router.POST("/api/sync/create", routes.TrackAppSyncItems)
	e.Router.POST("/api/focus/create", routes.TrackFocus)
	// e.Router.POST("/sync/track-items", routes.TrackAppItems)

	// Calendar
	e.Router.GET("/auth/calendar/redirect", func(e *core.RequestEvent) error {
		return routes.CalendarAuthHandler(app.CalendarService, e)
	})
	e.Router.POST("/auth/calendar/callback", func(e *core.RequestEvent) error {
		return routes.CalendarAuthCallback(app.CalendarService, e)
	})
	e.Router.POST("/api/calendar/sync", func(e *core.RequestEvent) error {
		return routes.CalendarSyncHandler(app.CalendarService, e)
	})

	// Mail
	e.Router.GET("/auth/mail/redirect", func(e *core.RequestEvent) error {
		return routes.MailAuthHandler(app.MailService, e)
	})

	e.Router.POST("/auth/mail/callback", func(e *core.RequestEvent) error {
		return routes.MailAuthCallback(app.MailService, e)
	})

	e.Router.POST("/api/mail/sync", func(e *core.RequestEvent) error {
		return routes.MailSyncHandler(app.MailService, e)
	})

	e.Router.GET("/api/mail/sync/status", func(e *core.RequestEvent) error {
		return routes.MailSyncStatusHandler(app.MailService, e)
	})

	// Money
	e.Router.POST("/api/fold/getotp", func(e *core.RequestEvent) error {
		return routes.FoldGetOtpHandler(app.FoldService, e)
	})
	e.Router.POST("/api/fold/verifyotp", func(e *core.RequestEvent) error {
		return routes.FoldVerifyOtpHandler(app.FoldService, e)
	})
	e.Router.GET("/api/fold/refresh", func(e *core.RequestEvent) error {
		return routes.FoldRefreshTokenHandler(app.FoldService, e)
	})
	e.Router.GET("/api/fold/user", func(e *core.RequestEvent) error {
		return routes.FoldUserHandler(app.FoldService, e)
	})
}

func (app *Application) InitCronjobs() error {
	cronJobs := []cronjobs.CronJob{
		{
			Name:     "track-device",
			Interval: "*/1 * * * *",
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
	app.Pb.RootCmd.AddCommand(cmd.NewSuperuserCommand(app.Pb))
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
