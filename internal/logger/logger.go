package logger

import (
	"fmt"
	"log"
	"log/slog"
	"os"

	"github.com/pocketbase/pocketbase"
)

const (
	flags      = log.Ldate | log.Ltime | log.Lshortfile
	boldRed    = "\033[1;31m"
	boldGreen  = "\033[1;32m"
	boldYellow = "\033[1;33m"
	boldPurple = "\033[1;35m"
	boldCyan   = "\033[1;36m"
	intenseRed = "\033[0;91m"
	reset      = "\033[0m"
)

var appLogger *slog.Logger

var (
	Debug   = log.New(os.Stdout, fmt.Sprintf("%s[DEBUG]%s ", boldCyan, reset), flags)
	Info    = log.New(os.Stdout, fmt.Sprintf("%s[INFO]%s ", boldGreen, reset), flags)
	Warning = log.New(os.Stdout, fmt.Sprintf("%s[WARNING]%s ", boldYellow, reset), flags)
	Error   = log.New(os.Stderr, fmt.Sprintf("%s[ERROR]%s ", boldRed, reset), flags)
	Fatal   = log.New(os.Stderr, fmt.Sprintf("%s[FATAL]%s", intenseRed, reset), flags)
)

func LogError(log string, message ...interface{}) {
	appLogger.Error(log, message...)
	fullMessage := []interface{}{log}
	fullMessage = append(fullMessage, message...)
	Error.Println(fullMessage...)
}

func LogWarning(log string, message ...interface{}) {
	appLogger.Warn(log, message...)
	fullMessage := []interface{}{log}
	fullMessage = append(fullMessage, message...)
	Warning.Println(fullMessage...)
}

func LogInfo(log string, message ...interface{}) {
	appLogger.Info(log, message...)
	fullMessage := []interface{}{log}
	fullMessage = append(fullMessage, message...)
	Info.Println(fullMessage...)
}

func RegisterApp(app *pocketbase.PocketBase) {
	appLogger = app.Logger()
}

func init() {

	stdout := os.Stdout
	stderr := os.Stderr

	Debug.SetOutput(stdout)
	Info.SetOutput(stdout)
	Error.SetOutput(stderr)
	Fatal.SetOutput(stderr)

	Debug.Println("Initialized Logger")
}
