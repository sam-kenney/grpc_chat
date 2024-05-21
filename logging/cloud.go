package logging

import (
	"context"
	"errors"
	"os"

	"cloud.google.com/go/logging"
)

// Google Cloud Logging wrapper with convenience methods for Debug through Error.
type CloudLogger struct {
	logger logging.Logger
	client *logging.Client
}

// Create a new Cloud CloudLogger.
// Requires "GCP_PROJECT" to be set in the environment.
func NewCloudLoggerFromEnv(ctx context.Context, name string) (*CloudLogger, error) {
	project := os.Getenv("GCP_PROJECT")

	if project == "" {
		return nil, errors.New("GCP_PROJECT not set in environment")
	}

	client, err := logging.NewClient(ctx, project)
	if err != nil {
		return nil, err
	}

	return &CloudLogger{
		logger: *client.Logger(name),
		client: client,
	}, nil
}

// Close all logging resources.
func (c CloudLogger) Close() error {
	return c.client.Close()
}

// Log an entry to Google Cloud Logging.
func (c CloudLogger) Log(entry Entry) {
	c.logger.Log(entry)
}

func (c CloudLogger) Debug(payload interface{}) {
	c.Log(logging.Entry{
		Severity: logging.Debug,
		Payload:  payload,
	})
}

func (c CloudLogger) Info(payload interface{}) {
	c.Log(logging.Entry{
		Severity: logging.Info,
		Payload:  payload,
	})
}

func (c CloudLogger) Warning(payload interface{}) {
	c.Log(logging.Entry{
		Severity: logging.Warning,
		Payload:  payload,
	})
}

func (c CloudLogger) Error(payload interface{}) {
	c.Log(logging.Entry{
		Severity: logging.Error,
		Payload:  payload,
	})
}

func (c CloudLogger) Critical(payload interface{}) {
	c.Log(logging.Entry{
		Severity: logging.Critical,
		Payload:  payload,
	})
}

// Raise a Critical error from a go error.
func (c CloudLogger) Err(err error) {
	c.Log(logging.Entry{
		Severity: logging.Critical,
		Payload:  err.Error(),
	})
}
