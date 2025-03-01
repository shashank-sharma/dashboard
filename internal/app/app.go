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
	"github.com/shashank-sharma/backend/internal/metrics"
	"github.com/shashank-sharma/backend/internal/routes"
	"github.com/shashank-sharma/backend/internal/services"
	"github.com/shashank-sharma/backend/internal/services/ai"
	"github.com/shashank-sharma/backend/internal/services/calendar"
	"github.com/shashank-sharma/backend/internal/services/feed"
	"github.com/shashank-sharma/backend/internal/services/fold"
	"github.com/shashank-sharma/backend/internal/services/mail"
	"github.com/shashank-sharma/backend/internal/services/providers"
	"github.com/shashank-sharma/backend/internal/services/workflow"
	"github.com/shashank-sharma/backend/internal/store"
)

type Application struct {
	Server          *http.Server
	Pb              *pocketbase.PocketBase
	FoldService     *fold.FoldService
	CalendarService *calendar.CalendarService
	MailService     *mail.MailService
	WorkflowEngine  *workflow.WorkflowEngine
	FeedService     *services.FeedService
}

func defaultPublicDir() string {
	if strings.HasPrefix(os.Args[0], os.TempDir()) {
		return "./pb_public"
	}

	return filepath.Join(os.Args[0], "../pb_public")
}

func New() *Application {
	pb := pocketbase.New()
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Failed to load environment variables")
	}

	// Define all the service here
	foldService := fold.NewFoldService("https://api.fold.money/api")
	calendarService := calendar.NewCalendarService()
	mailService := mail.NewMailService()
	workflowEngine := workflow.NewWorkflowEngine(pb)

	aiConfig := config.GetAIConfig()
	var aiClient ai.AIClient
	
	if aiConfig.Service != config.AIServiceNone {
		var err error
		aiClient, err = ai.NewAIClient(aiConfig)
		if err != nil {
			logger.LogError("Failed to initialize AI client: " + err.Error())
			logger.LogInfo("Continuing without AI functionality")
		}
		logger.LogInfo("AI client initialized")
	}
	
	// Initialize the feed processor with AI client
	processor := feed.NewFeedProcessor(aiClient)
	feedService := feed.NewFeedService(processor)

	// Register feed source providers
	feedService.RegisterProvider(providers.NewRSSProvider())
	feedService.RegisterProvider(providers.NewHackerNewsProvider())

	app := &Application{
		Pb:              pb,
		FoldService:     foldService,
		CalendarService: calendarService,
		MailService:     mailService,
		WorkflowEngine:  workflowEngine,
		FeedService: &feedService,
	}

	pb.OnServe().BindFunc(func(e *core.ServeEvent) error {
		// Logger need to be initialized inside pocketbase
		// before serve
		logger.RegisterApp(pb)
		store.InitApp(pb)
		config.Init(app.Pb)
		metrics.RegisterPrometheusMetrics(pb)		
		app.InitCronjobs()
		app.configureRoutes(e)

		if err := metrics.StartMetricsServer(app.Pb); err != nil {
			logger.LogError("Failed to start metrics server", "error", err)
		}

		return e.Next()
	})

	app.registerHooks()
	return app
}

func (app *Application) configureRoutes(e *core.ServeEvent) {
	routes.RegisterWorkflowRoutes(e, app.WorkflowEngine)
	routes.RegisterFeedRoutes(e, *app.FeedService)
	routes.RegisterCredentialRoutes(e)
	routes.RegisterTrackRoutes(e)
	routes.RegisterCalendarRoutes(e, app.CalendarService)
	routes.RegisterMailRoutes(e, app.MailService)
	routes.RegisterFoldRoutes(e, app.FoldService)
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
	app.Pb.RootCmd.PersistentFlags().BoolVar(
		&config.EnableMetricsFlag,
		"metrics",
		false,
		"enable Prometheus metrics collection",
	)

	isGoRun := strings.HasPrefix(os.Args[0], os.TempDir())
	logger.Debug.Println("isGoRun:", isGoRun)

	migratecmd.MustRegister(app.Pb, app.Pb.RootCmd, migratecmd.Config{
		// enable auto creation of migration files when making collection changes in the Admin UI
		// (the isGoRun check is to enable it only during development)
		Automigrate: isGoRun,
	})

	return app.Pb.Execute()
}
