package app

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/joho/godotenv"
	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/apis"
	"github.com/pocketbase/pocketbase/core"
	"github.com/shashank-sharma/backend/internal/config"
	"github.com/shashank-sharma/backend/internal/cronjobs"
	"github.com/shashank-sharma/backend/internal/gui"
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
	postInitHooks   []func()
}

func New() *Application {
	pb := pocketbase.NewWithConfig(pocketbase.Config{
		DefaultDataDir: "./pb_data",
		HideStartBanner: false,
		DefaultDev: false,
	})
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Failed to load environment variables")
	}

	// Initialize store and config (basic initialization only)
	store.InitApp(pb)
	config.Init(pb)

	// Create a minimal Application with just PocketBase
	app := &Application{
		Pb:              pb,
		postInitHooks:   make([]func(), 0),
	}

	// Experiment: Register a post-initialization hook
	app.AddPostInitHook(func() {
		logger.LogInfo("This is an example post-initialization hook - app is fully ready")
	})

	pb.OnServe().BindFunc(func(e *core.ServeEvent) error {
		// STAGE 1: Initialize base services that don't depend on application services
		logger.InitLog(pb)
		logger.LogInfo("Base logging initialized")
		
		metrics.RegisterPrometheusMetrics(pb)
		logger.LogInfo("Metrics initialized")
		
		logger.LogInfo("Initializing application services...")
			
		app.initializeServices()		
		app.InitCronjobs()
		app.configureRoutes(e)
		
		if err := metrics.StartMetricsServer(app.Pb); err != nil {
			logger.LogError("Failed to start metrics server", "error", err)
		}
		
		logger.LogInfo("All application services initialized")
		app.RunPostInitHooks()

		return e.Next()
	})

	app.registerHooks()
	return app
}

// initializeServices creates and initializes all application services
func (app *Application) initializeServices() {
	app.FoldService = fold.NewFoldService("https://api.fold.money/api")
	app.CalendarService = calendar.NewCalendarService()
	app.MailService = mail.NewMailService()
	app.WorkflowEngine = workflow.NewWorkflowEngine(app.Pb)

	aiConfig := config.GetAIConfig()
	var aiClient ai.AIClient
	
	if aiConfig.Service != config.AIServiceNone {
		var err error
		aiClient, err = ai.NewAIClient(aiConfig)
		if err != nil {
			logger.LogError("Failed to initialize AI client: " + err.Error())
			logger.LogInfo("Continuing without AI functionality")
		} else {
			logger.LogInfo("AI client initialized")
		}
	}
	
	// Initialize the feed processor with AI client
	processor := feed.NewFeedProcessor(aiClient)
	feedService := feed.NewFeedService(processor)
	feedService.RegisterProvider(providers.NewRSSProvider())
	feedService.RegisterProvider(providers.NewHackerNewsProvider())
	
	app.FeedService = &feedService
	
	logger.LogInfo("All services initialized successfully")
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
	// Check if GUI flag is set
	withGUI, ok := app.Pb.Store().Get("WITH_GUI").(bool)
	fileLogging, okLogging := app.Pb.Store().Get("FILE_LOGGING_ENABLED").(bool)
	
	if ok && withGUI && okLogging && fileLogging {
		logFilePath, _ := app.Pb.Store().Get("LOG_FILE_PATH").(string)
		
		// Start the server in a goroutine
		go func() {
			app.Pb.Bootstrap()
			err := app.Serve()

			if err != nil {
				logger.LogInfo("Server closed error: " + err.Error())
			}
		}()
		
		// Allow some time for the server to start
		time.Sleep(500 * time.Millisecond)

		guiStatus := gui.GUIStatus{
			ServerRunning: true,
			MetricsEnabled: app.Pb.Store().Get("METRICS_ENABLED").(bool),
		}
		
		metadata := app.collectServerMetadata()
		return gui.StartGUI(logFilePath, guiStatus, metadata)
	}
	
	// Default behavior (no GUI)
	return app.Serve()
}

// collectServerMetadata gathers information about the server for display in the GUI
func (app *Application) collectServerMetadata() gui.ServerMetadata {
	// Collect basic server info
	serverURL := "http://localhost:8090"
	if customURL, ok := app.Pb.Store().Get("SERVER_URL").(string); ok && customURL != "" {
		serverURL = customURL
	}
	
	// Get environment variables	
	// Get all keys from the store
	envVars := app.Pb.Store().GetAll()
	
	// Get current environment
	environment := "production"
	if env, ok := app.Pb.Store().Get("APP_ENVIRONMENT").(string); ok {
		environment = env
	}
	
	// Collect configured cron jobs
	cronJobs := []gui.CronJob{}
	for _, job := range cronjobs.GetActiveJobs() {
		cronJobs = append(cronJobs, gui.CronJob{
			Name:     job.Name,
			Schedule: job.Interval,
			LastRun:  job.LastRun,
			Status:   job.GetStatusString(),
		})
	}
	
	// Collect API endpoints
	endpoints := []string{
		"/api/collections",
		"/api/admins",
		"/api/feeds",
		"/api/workflows",
	}
	
	// Create the metadata struct
	metadata := gui.ServerMetadata{
		ServerURL:      serverURL,
		ServerVersion:  "v1.0.0", // You can retrieve this from a version package if available
		Environment:    environment,
		EnvVariables:   envVars,
		CronJobs:       cronJobs,
		StartTime:      time.Now(), // This should ideally be the actual server start time
		DataDirectory:  "./pb_data",
		APIEndpoints:   endpoints,
	}
	
	return metadata
}

func (app *Application) Serve() error {
	app.Pb.Bootstrap()
	err := apis.Serve(app.Pb, apis.ServeConfig{
		HttpAddr:           "0.0.0.0:8090",
		ShowStartBanner:    false,
	})

	if errors.Is(err, http.ErrServerClosed) {
		return nil
	}

	return err
}

// AddPostInitHook adds a function to be executed after the server is fully initialized
func (app *Application) AddPostInitHook(hookFunc func()) {
	app.postInitHooks = append(app.postInitHooks, hookFunc)
}

// RunPostInitHooks executes all registered post-initialization hooks
func (app *Application) RunPostInitHooks() {
	logger.LogInfo("Running post-initialization hooks")
	for _, hook := range app.postInitHooks {
		hook()
	}
}