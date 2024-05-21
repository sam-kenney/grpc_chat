package logging

import (
	"context"
	"fmt"
	"os"
	"strings"
)

// Create a new `logging.StandardLogger`.
// If the `NO_COLOR` environment variable is set, no colours will be output.
func NewStandardLogger(ctx context.Context, severity Severity, name string) StandardLogger {
	useColour := true
	colour := os.Getenv("NO_COLOR")
	if colour != "" {
		useColour = false
	}

	return StandardLogger{ctx, severity, name, useColour}
}

type StandardLogger struct {
	ctx       context.Context
	Severity  Severity
	Name      string
	useColour bool
}

func (l StandardLogger) Log(severity Severity, payload interface{}) {
	// Don't show messages for severities less than configured
	if severity < l.Severity {
		return
	}

	sev := strings.ToUpper(severity.String())[:4]
	colour, ok := colours[severity]

	if !ok || !l.useColour {
		fmt.Printf("[%s] %s: %s\n", sev, l.Name, payload)
		return
	}

	colour.Printf("[%s]", sev)
	fmt.Printf(" %s: %s\n", l.Name, payload)
}

func (l StandardLogger) Debug(payload interface{}) {
	l.Log(Debug, payload)
}

func (l StandardLogger) Info(payload interface{}) {
	l.Log(Info, payload)
}

func (l StandardLogger) Warning(payload interface{}) {
	l.Log(Warning, payload)
}

func (l StandardLogger) Error(payload interface{}) {
	l.Log(Error, payload)
}

func (l StandardLogger) Critical(payload interface{}) {
	l.Log(Critical, payload)
}

func (l StandardLogger) Err(err error) {
	l.Log(Critical, err.Error())
}
