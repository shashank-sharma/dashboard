package logger

import (
	"fmt"
	"log"
	"log/slog"
	"os"
	"runtime"
	"strings"
	"sync"

	"github.com/pocketbase/pocketbase"
)

// ANSI color codes
const (
	boldRed    = "\033[1;31m"
	boldGreen  = "\033[1;32m"
	boldYellow = "\033[1;33m"
	boldCyan   = "\033[1;36m"
	intenseRed = "\033[0;91m"
	reset      = "\033[0m"
)

// Logger flags - we only want date and time, as we handle caller info manually
const callerFlags = log.Ldate | log.Ltime

var (
	appLogger *slog.Logger
	// Initialize standard loggers
	Debug   = log.New(os.Stdout, fmt.Sprintf("%s[DEBUG]%s ", boldCyan, reset), callerFlags)
	Info    = log.New(os.Stdout, fmt.Sprintf("%s[INFO]%s ", boldGreen, reset), callerFlags)
	Warning = log.New(os.Stdout, fmt.Sprintf("%s[WARNING]%s ", boldYellow, reset), callerFlags)
	Error   = log.New(os.Stderr, fmt.Sprintf("%s[ERROR]%s ", boldRed, reset), callerFlags)
	Fatal   = log.New(os.Stderr, fmt.Sprintf("%s[FATAL]%s ", intenseRed, reset), callerFlags)
	// Pre-allocate buffer for formatKeyValuePairs
	bufPool = sync.Pool{
		New: func() interface{} {
			return new(strings.Builder)
		},
	}
)

// getCallerInfo returns file:line from the caller's perspective
func getCallerInfo(skip int) string {
	_, file, line, ok := runtime.Caller(skip + 1)
	if !ok {
		return "???:0"
	}
	// Extract just the filename for cleaner logs
	short := file
	for i := len(file) - 1; i > 0; i-- {
		if file[i] == '/' {
			short = file[i+1:]
			break
		}
	}
	return fmt.Sprintf("%s:%d", short, line)
}

// formatKeyValuePairs formats message arguments as key-value pairs
// If the arguments don't follow key-string/value pairs, they'll be marked as badKey
func formatKeyValuePairs(args ...interface{}) string {
	if len(args) == 0 {
		return ""
	}

	// Get a buffer from the pool
	builder := bufPool.Get().(*strings.Builder)
	builder.Reset()
	defer bufPool.Put(builder)

	for i := 0; i < len(args); i += 2 {
		// Check if we have a valid key-value pair
		if i+1 >= len(args) {
			// Odd number of arguments
			fmt.Fprintf(builder, "{badKey: %v} ", args[i])
			continue
		}

		key, ok := args[i].(string)
		if !ok {
			// Key is not a string
			fmt.Fprintf(builder, "{badKey: %v, value: %v} ", args[i], args[i+1])
			continue
		}

		// Format as key-value pair
		fmt.Fprintf(builder, "{%s: %v} ", key, args[i+1])
	}

	return builder.String()
}

// LogDebug logs with the DEBUG level and correct caller information
// Handles key-value pairs in the format: LogInfo("message", "key1", value1, "key2", value2)
func LogDebug(message string, args ...interface{}) {
	if appLogger != nil {
		appLogger.Debug(message, args...)
	}

	callerInfo := getCallerInfo(1)
	if len(args) > 0 {
		Debug.Println(callerInfo + ": " + message + " " + formatKeyValuePairs(args...))
	} else {
		Debug.Println(callerInfo + ": " + message)
	}
}

// LogInfo logs with the INFO level and correct caller information
// Handles key-value pairs in the format: LogInfo("message", "key1", value1, "key2", value2)
func LogInfo(message string, args ...interface{}) {
	if appLogger != nil {
		appLogger.Info(message, args...)
	}

	callerInfo := getCallerInfo(1)
	if len(args) > 0 {
		Info.Println(callerInfo + ": " + message + " " + formatKeyValuePairs(args...))
	} else {
		Info.Println(callerInfo + ": " + message)
	}
}

// LogError logs with the ERROR level and correct caller information
// Handles key-value pairs in the format: LogError("message", "key1", value1, "key2", value2)
func LogError(message string, args ...interface{}) {
	if appLogger != nil {
		appLogger.Error(message, args...)
	}

	callerInfo := getCallerInfo(1)
	if len(args) > 0 {
		Error.Println(callerInfo + ": " + message + " " + formatKeyValuePairs(args...))
	} else {
		Error.Println(callerInfo + ": " + message)
	}
}

// LogWarning logs with the WARNING level and correct caller information
// Handles key-value pairs in the format: LogWarning("message", "key1", value1, "key2", value2)
func LogWarning(message string, args ...interface{}) {
	if appLogger != nil {
		appLogger.Warn(message, args...)
	}

	callerInfo := getCallerInfo(1)
	if len(args) > 0 {
		Warning.Println(callerInfo + ": " + message + " " + formatKeyValuePairs(args...))
	} else {
		Warning.Println(callerInfo + ": " + message)
	}
}

// RegisterApp registers the pocketbase app for logging
func RegisterApp(app *pocketbase.PocketBase) {
	appLogger = app.Logger()
}

func init() {
	Debug.SetOutput(os.Stdout)
	Info.SetOutput(os.Stdout)
	Warning.SetOutput(os.Stdout)
	Error.SetOutput(os.Stderr)
	Fatal.SetOutput(os.Stderr)

	Debug.Println(getCallerInfo(0) + ": Initialized Logger")
}
