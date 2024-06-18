package logger

import (
	"fmt"
	"log/slog"

	"github.com/pocketbase/pocketbase"
)

/*
const (
	flags = log.Ldate | log.Ltime | log.Lshortfile
)
*/

var appLogger *slog.Logger

/*
var (
	Debug     = log.New(os.Stdout, "[DEBUG] ", flags)
	Info      = log.New(os.Stdout, "[INFO] ", flags)
	Warning   = log.New(os.Stdout, "[WARNING] ", flags)
	Error     = log.New(os.Stderr, "[ERROR] ", flags)
	Fatal     = log.New(os.Stderr, "[FATAL]", flags)
)
*/

func LogError(message ...interface{}) {
	fmt.Println("Logger: ", appLogger)
	appLogger.Error("", message...)
	// Error.Println(message...)
}

func RegisterApp(app *pocketbase.PocketBase) {
	appLogger = app.Logger()
}

func init() {

	/*
		stdout := os.Stdout
		stderr := os.Stderr

		Debug.SetOutput(stdout)
		Info.SetOutput(stdout)
		Error.SetOutput(stderr)
		Fatal.SetOutput(stderr)

		log.SetOutput(Debug.Writer())
		log.SetPrefix("[DEBUG]")
		log.SetFlags(flags)
	*/

	// Debug.Println("Initialized Logger")
}
