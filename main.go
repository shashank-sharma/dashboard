package main

import (
	"log"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"bufio"
	"github.com/labstack/echo/v5"
	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/apis"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/plugins/jsvm"
	"github.com/pocketbase/pocketbase/plugins/migratecmd"
)


func streamMP3(path string) (io.Reader, error) {
    f, err := os.Open(path)
    if err != nil {
        return nil, err
    }
    return f, nil
}

func streamMP3Handler(c echo.Context) error {
    path := c.QueryParam("path")
    if path == "" {
        return c.JSON(http.StatusBadRequest, map[string]string{"error": "path parameter is required"})
    }

    reader, err := streamMP3(path)
    if err != nil {
        return c.JSON(http.StatusInternalServerError, map[string]string{"error": fmt.Sprintf("error streaming mp3: %v", err)})
    }

    // Create a buffered reader with a buffer size of 4096 bytes.
    bufferedReader := bufio.NewReaderSize(reader, 4096)

    c.Response().Header().Set("Content-Type", "audio/mpeg")
    _, err = io.CopyBuffer(c.Response(), bufferedReader, nil)
    if err != nil {
        return err
    }
    return nil
}


func defaultPublicDir() string {
	if strings.HasPrefix(os.Args[0], os.TempDir()) {
		// most likely ran with go run
		return "./pb_public"
	}

	return filepath.Join(os.Args[0], "../pb_public")
}

func main() {
	app := pocketbase.New()

	var publicDirFlag string

	// add "--publicDir" option flag
	app.RootCmd.PersistentFlags().StringVar(
		&publicDirFlag,
		"publicDir",
		defaultPublicDir(),
		"the directory to serve static files",
	)
	migrationsDir := "" // default to "pb_migrations" (for js) and "migrations" (for go)

	// load js files to allow loading external JavaScript migrations
	jsvm.MustRegisterMigrations(app, &jsvm.MigrationsOptions{
		Dir: migrationsDir,
	})

	// register the `migrate` command
	migratecmd.MustRegister(app, app.RootCmd, &migratecmd.Options{
		TemplateLang: migratecmd.TemplateLangJS, // or migratecmd.TemplateLangGo (default)
		Dir:          migrationsDir,
		Automigrate:  true,
	})

	app.OnBeforeServe().Add(func(e *core.ServeEvent) error {
		// serves static files from the provided public dir (if exists)
		e.Router.GET("/*", apis.StaticDirectoryHandler(os.DirFS(publicDirFlag), true))

		e.Router.AddRoute(echo.Route{
			Method: http.MethodGet,
			Path:   "/api/hello",
			Handler: func(c echo.Context) error {
				obj := map[string]interface{}{"message": "Hello world!"}
				return c.JSON(http.StatusOK, obj)
			},
			// Middlewares: []echo.MiddlewareFunc{
			// 	apis.RequireAdminOrUserAuth(),
			// },
		})


		    e.Router.AddRoute(echo.Route{
        		Method: http.MethodGet,
        		Path:  "/stream_mp3",
        		Handler: streamMP3Handler,
    		})

		return nil
	})

	if err := app.Start(); err != nil {
		log.Fatal(err)
	}
}
