package util

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/pocketbase/pocketbase"
)

var app *pocketbase.PocketBase

type Config struct {
	// PocketBase server URL
	PocketBaseURL string `envconfig:"POCKETBASE_URL"`

	// PocketBase Admin Key
	PocketBaseAdminKey string `envconfig:"POCKETBASE_ADMIN_KEY"`

	// PocketBase Database Name
	PocketBaseDatabaseName string `envconfig:"POCKETBASE_DATABASE_NAME"`

	// Application Port
	Port int `envconfig:"PORT"`
}

func defaultPublicDir() string {
	if strings.HasPrefix(os.Args[0], os.TempDir()) {
		// most likely ran with go run
		return "./pb_public"
	}

	return filepath.Join(os.Args[0], "../pb_public")
}

func Init() {
	app = pocketbase.New()

	var publicDirFlag string

	// add "--publicDir" option flag
	app.RootCmd.PersistentFlags().StringVar(
		&publicDirFlag,
		"publicDir",
		defaultPublicDir(),
		"the directory to serve static files",
	)

	fmt.Println("Created add: ", app)
}

func GetApp() *pocketbase.PocketBase {
	fmt.Println("Found app: ", app)
	return app
}
