package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"github.com/pocketbase/pocketbase"
)

var EnableMetricsFlag bool
var FileLoggingFlag bool
var WithGuiFlag bool

type ConfigFlags struct {
	Metrics     bool
	FileLogging bool
	WithGui     bool
	Dev         bool
	HttpAddr    string
}

func getEnv(key string, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

func Init(pb *pocketbase.PocketBase, config ConfigFlags) {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Failed to load environment variables")
	}

	pb.Store().Set("ENCRYPTION_KEY", getEnv("ENCRYPTION_KEY", "default_encryption_key"))
	pb.Store().Set("METRICS_ENABLED", config.Metrics)
	pb.Store().Set("METRICS_PORT", getEnv("METRICS_PORT", "9091"))
	pb.Store().Set("FILE_LOGGING_ENABLED", config.FileLogging)
	pb.Store().Set("LOG_FILE_PATH", getEnv("LOG_FILE_PATH", "logs/app.log"))
	pb.Store().Set("WITH_GUI", config.WithGui)
	pb.Store().Set("DEV", config.Dev)
	pb.Store().Set("HTTP_ADDR", config.HttpAddr)

	// Set global flags for easy access
	EnableMetricsFlag = config.Metrics
	FileLoggingFlag = config.FileLogging
	WithGuiFlag = config.WithGui
}
