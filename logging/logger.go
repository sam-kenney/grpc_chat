package logging

import (
	"context"
	"os"

	"cloud.google.com/go/logging"
	"github.com/fatih/color"
)

// Re-exports stuff from google cloud logging
// so end user doesn't have to import two logging libs.
type Entry = logging.Entry
type Severity = logging.Severity

const (
	Default   = logging.Default
	Debug     = logging.Debug
	Info      = logging.Info
	Warning   = logging.Warning
	Error     = logging.Error
	Critical  = logging.Critical
	Alert     = logging.Alert
	Emergency = logging.Emergency
)

// Map colours to Severity levels.
var colours = map[Severity]*color.Color{
	Default:   color.New(color.FgWhite),
	Debug:     color.New(color.FgWhite),
	Info:      color.New(color.FgCyan),
	Warning:   color.New(color.FgYellow),
	Error:     color.New(color.FgRed),
	Critical:  color.New(color.FgRed),
	Alert:     color.New(color.FgRed),
	Emergency: color.New(color.FgRed),
}

// Interface to derive logging clients from.
type Logger interface {
	Debug(interface{})
	Info(interface{})
	Warning(interface{})
	Error(interface{})
	Critical(interface{})
	Err(error)
}

// Create a new logger from the environment.
// Uses the `LOGGING_ENVIRONMENT` environment variable to determine
// how to log messages. If set to "GCP", a `logging.CloudLogger` is used,
// otherwise, a `logging.StandardLogger` is used.
func NewLoggerFromEnv(ctx context.Context, severity Severity, name string) (Logger, error) {
	env := os.Getenv("LOGGING_ENVIRONMENT")
	if env == "GCP" {
		return NewCloudLoggerFromEnv(ctx, name)
	}

	return NewStandardLogger(ctx, severity, name), nil
}
