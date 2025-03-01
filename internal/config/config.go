package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"github.com/pocketbase/pocketbase"
)

var EnableMetricsFlag bool

func getEnv(key string, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

func Init(pb *pocketbase.PocketBase) {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Failed to load environment variables")
	}

	pb.Store().Set("ENCRYPTION_KEY", getEnv("ENCRYPTION_KEY", "default_encryption_key"))
	pb.Store().Set("METRICS_ENABLED", EnableMetricsFlag)
	pb.Store().Set("METRICS_PORT", getEnv("METRICS_PORT", "9091"))
}
