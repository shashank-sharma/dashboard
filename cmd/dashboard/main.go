package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/pocketbase/pocketbase/plugins/migratecmd"
	"github.com/shashank-sharma/backend/internal/app"
	"github.com/shashank-sharma/backend/internal/config"
	_ "github.com/shashank-sharma/backend/migrations"
)

func main() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [OPTIONS]\n", os.Args[1])
		fmt.Fprintf(os.Stderr, "\n")
		fmt.Fprintf(os.Stderr, "This application is a backend service with configurable options.\n")
		fmt.Fprintf(os.Stderr, "\n")
		fmt.Fprintf(os.Stderr, "Options:\n")
		flag.PrintDefaults()
		fmt.Fprintf(os.Stderr, "\n")
		fmt.Fprintf(os.Stderr, "Examples:\n")
		fmt.Fprintf(os.Stderr, "  Run in production mode:  %s\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  Run in development mode: %s -dev\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  Run with GUI and metrics: %s -with-gui -metrics\n", os.Args[0])
	}

    metrics := flag.Bool("metrics", false, "Enable metrics")
    fileLogging := flag.Bool("file-logging", false, "Enable file logging")
    withGui := flag.Bool("with-gui", false, "Enable GUI")
    dev := flag.Bool("dev", false, "Run in development mode")

    flag.Parse()

	args := flag.Args()
	config := config.ConfigFlags{
		Metrics:     *metrics,
		FileLogging: *fileLogging,
		WithGui:     *withGui,
		Dev:         *dev,
	}
	application := app.New(config)


	migratecmd.MustRegister(application.Pb, application.Pb.RootCmd, migratecmd.Config{
		// enable auto creation of migration files when making collection changes in the Dashboard
		// (the isGoRun check is to enable it only during development)
		Automigrate: strings.HasPrefix(os.Args[0], os.TempDir()),
	})
    if len(args) > 0 && args[0] == "migrate" {
        fmt.Println("Running migration...")
		application.Pb.Start()
        // You can add your migration logic here or call a migration function
        return
    } else {
		application.Start()
	}
}
