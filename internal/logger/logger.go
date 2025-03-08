package logger

import (
	"fmt"
	"log"
	"log/slog"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/fatih/color"
	"github.com/pocketbase/pocketbase"
)

const callerFlags = log.Ldate | log.Ltime

var (
	appLogger *slog.Logger
	isDev     bool
	
	Debug   = log.New(os.Stdout, "", 0)
	Info    = log.New(os.Stdout, "", 0)
	Warning = log.New(os.Stdout, "", 0)
	Error   = log.New(os.Stderr, "", 0)
	Fatal   = log.New(os.Stderr, "", 0)
	
	dateTimeColor = color.New(color.FgHiBlack)
	
	// Pre-allocate buffer for formatKeyValuePairs
	bufPool = sync.Pool{
		New: func() interface{} {
			return new(strings.Builder)
		},
	}

	logFile *os.File

	// Initialize fatih/color styles for log levels
	debugLabel   = color.New(color.FgHiBlack, color.Bold)
	debugText    = color.New(color.FgHiBlack)
	infoLabel    = color.New(color.FgGreen, color.Bold)
	infoText     = color.New(color.Reset)
	warningLabel = color.New(color.FgYellow, color.Bold)
	warningText  = color.New(color.FgYellow)
	errorLabel   = color.New(color.FgRed, color.Bold)
	errorText    = color.New(color.FgRed)
	fatalLabel   = color.New(color.FgHiRed, color.Bold)
	fatalText    = color.New(color.FgHiRed)
	
	// Color styles for key-value pairs
	debugKey     = color.New(color.FgHiBlack, color.Bold)
	debugValue   = color.New(color.FgHiBlack)
	warningKey   = color.New(color.FgYellow, color.Bold)
	warningValue = color.New(color.FgYellow)
	errorKey     = color.New(color.FgRed, color.Bold)
	errorValue   = color.New(color.FgRed)
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

func formatKeyValuePairs(logLevel string, args ...interface{}) string {
	if len(args) == 0 {
		return ""
	}

	// Get a buffer from the pool
	builder := bufPool.Get().(*strings.Builder)
	builder.Reset()
	defer bufPool.Put(builder)

	// Add a newline before the first pair for better readability
	builder.WriteString("\n")

	pairCount := len(args) / 2
	if len(args)%2 != 0 {
		pairCount++
	}

	var keyFormatter, valueFormatter *color.Color
	switch logLevel {
	case "ERROR":
		keyFormatter = errorKey
		valueFormatter = errorValue
	case "FATAL":
		keyFormatter = errorKey
		valueFormatter = errorValue
	case "WARNING":
		keyFormatter = warningKey
		valueFormatter = warningValue
	case "INFO":
		keyFormatter = debugKey
		valueFormatter = debugValue
	case "DEBUG":
		keyFormatter = debugKey
		valueFormatter = debugValue
	default:
		keyFormatter = debugKey
		valueFormatter = debugValue
	}

	for i := 0; i < len(args); i += 2 {
		currentPair := i / 2
		treeSymbol := "├─ "
		if currentPair == pairCount-1 {
			treeSymbol = "└─ "
		}

		builder.WriteString("  " + treeSymbol)

		// Check if we have a valid key-value pair
		if i+1 >= len(args) {
			// Odd number of arguments
			builder.WriteString(keyFormatter.Sprintf("badKey: ") + valueFormatter.Sprintf("%v\n", args[i]))
			continue
		}

		key, ok := args[i].(string)
		if !ok {
			// Key is not a string
			builder.WriteString(keyFormatter.Sprintf("badKey: ") + valueFormatter.Sprintf("%v, value: %v\n", args[i], args[i+1]))
			continue
		}

		builder.WriteString(keyFormatter.Sprintf("%s: ", key) + valueFormatter.Sprintf("%v\n", args[i+1]))
	}

	return builder.String()
}

// formatLogPrefix returns the formatted log prefix with appropriate colors
func formatLogPrefix(level string, callerInfo string) string {
	timestamp := dateTimeColor.Sprintf("%s ", time.Now().Format("2006/01/02 15:04:05"))
	switch level {
	case "DEBUG":
		return timestamp + debugLabel.Sprintf("[DEBUG] ") + debugText.Sprintf("%s: ", callerInfo)
	case "INFO":
		return timestamp + infoLabel.Sprintf("[INFO] ") + infoText.Sprintf("%s: ", callerInfo)
	case "WARNING":
		return timestamp + warningLabel.Sprintf("[WARNING] ") + warningText.Sprintf("%s: ", callerInfo)
	case "ERROR":
		return timestamp + errorLabel.Sprintf("[ERROR] ") + errorText.Sprintf("%s: ", callerInfo)
	case "FATAL":
		return timestamp + fatalLabel.Sprintf("[FATAL] ") + fatalText.Sprintf("%s: ", callerInfo)
	default:
		return timestamp +fmt.Sprintf("[%s] %s: ", level, callerInfo)
	}
}

// LogDebug logs with the DEBUG level and correct caller information
// Handles key-value pairs in the format: LogInfo("message", "key1", value1, "key2", value2)
func LogDebug(message string, args ...interface{}) {
	if appLogger != nil && isDev {
		appLogger.Debug(message, args...)
	}

	callerInfo := getCallerInfo(1)
	prefix := formatLogPrefix("DEBUG", callerInfo)
	
	if len(args) > 0 {
		Debug.Println(prefix + debugText.Sprint(message) + formatKeyValuePairs("DEBUG", args...))
	} else {
		Debug.Println(prefix + debugText.Sprint(message))
	}
}

// LogInfo logs with the INFO level and correct caller information
// Handles key-value pairs in the format: LogInfo("message", "key1", value1, "key2", value2)
func LogInfo(message string, args ...interface{}) {
	if appLogger != nil && isDev {
		appLogger.Info(message, args...)
	}

	callerInfo := getCallerInfo(1)
	prefix := formatLogPrefix("INFO", callerInfo)
	
	if len(args) > 0 {
		Info.Println(debugText.Sprint(prefix) + infoText.Sprint(message) + formatKeyValuePairs("INFO", args...))
	} else {
		Info.Println(debugText.Sprint(prefix) + infoText.Sprint(message))
	}
}

// LogError logs with the ERROR level and correct caller information
// Handles key-value pairs in the format: LogError("message", "key1", value1, "key2", value2)
func LogError(message string, args ...interface{}) {
	if appLogger != nil && isDev {
		appLogger.Error(message, args...)
	}

	callerInfo := getCallerInfo(1)
	prefix := formatLogPrefix("ERROR", callerInfo)
	
	if len(args) > 0 {
		Error.Println(prefix + errorText.Sprint(message) + formatKeyValuePairs("ERROR", args...))
	} else {
		Error.Println(prefix + errorText.Sprint(message))
	}
}

// LogWarning logs with the WARNING level and correct caller information
// Handles key-value pairs in the format: LogWarning("message", "key1", value1, "key2", value2)
func LogWarning(message string, args ...interface{}) {
	if appLogger != nil && isDev {
		appLogger.Warn(message, args...)
	}

	callerInfo := getCallerInfo(1)
	prefix := formatLogPrefix("WARNING", callerInfo)
	
	if len(args) > 0 {
		Warning.Println(prefix + warningText.Sprint(message) + formatKeyValuePairs("WARNING", args...))
	} else {
		Warning.Println(prefix + warningText.Sprint(message))
	}
}

// LogFatal logs with the FATAL level and correct caller information
// Handles key-value pairs in the format: LogFatal("message", "key1", value1, "key2", value2)
func LogFatal(message string, args ...interface{}) {
	if appLogger != nil && isDev {
		appLogger.Error(message, args...)
	}

	callerInfo := getCallerInfo(1)
	prefix := formatLogPrefix("FATAL", callerInfo)
	
	if len(args) > 0 {
		Fatal.Println(prefix + fatalText.Sprint(message) + formatKeyValuePairs("FATAL", args...))
	} else {
		Fatal.Println(prefix + fatalText.Sprint(message))
	}
	
	os.Exit(1)
}

// RegisterApp registers the pocketbase app for logging
func RegisterApp(app *pocketbase.PocketBase) {
	appLogger = app.Logger()
}

// Cleanup closes the log file if it's open
func Cleanup() {
	if logFile != nil {
		logFile.Close()
	}
}

func InitLog(app *pocketbase.PocketBase) {
	appLogger = app.Logger()
	
	// Check if file logging is enabled
	fileLoggingEnabled, _ := app.Store().Get("FILE_LOGGING_ENABLED").(bool)
	isDev, _ := app.Store().Get("DEV").(bool)	
	if fileLoggingEnabled {
		logFilePath, _ := app.Store().Get("LOG_FILE_PATH").(string)
		
		// Create directory if it doesn't exist
		dir := filepath.Dir(logFilePath)
		if err := os.MkdirAll(dir, 0755); err != nil {
			fmt.Fprintf(os.Stderr, "Failed to create log directory: %v\n", err)
			return
		}
		
		// Open log file
		file, err := os.OpenFile(logFilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to open log file: %v\n", err)
			return
		}
		
		// Store file for later cleanup
		logFile = file
		
		// Redirect loggers to file
		Debug.SetOutput(file)
		Info.SetOutput(file)
		Warning.SetOutput(file)
		Error.SetOutput(file)
		Fatal.SetOutput(file)
		
		LogInfo("File logging enabled, writing to " + logFilePath)
	} else {
		Debug.SetOutput(os.Stdout)
		Info.SetOutput(os.Stdout)
		Warning.SetOutput(os.Stdout)
		Error.SetOutput(os.Stderr)
		Fatal.SetOutput(os.Stderr)

		LogInfo("Development mode", "isDev", isDev)
		LogInfo("Console logging enabled")
	}
}