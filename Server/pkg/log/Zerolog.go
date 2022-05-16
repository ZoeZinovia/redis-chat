package log

import (
	"fmt"
	"os"
	"runtime"
	"strings"

	"github.com/rs/zerolog"
)

type ZeroLog struct {
	logger *zerolog.Logger
}

// Creates a new logger with the required
// configurations
func CreateLogger() Log {
	// Define log level
	zerolog.SetGlobalLevel(zerolog.InfoLevel)

	// Create logger with timestamp
	logger := zerolog.New(os.Stderr).With().Timestamp().Logger()

	// Set format of logger to pretty
	logger = logger.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	// Return ZeroLog instance
	return &ZeroLog{
		logger: &logger,
	}
}

func (l *ZeroLog) Info(msg string, information Information) {
	l.logger.Info().Str("location", getLocation(2)).Fields(information).Msg(msg)
}

func (l *ZeroLog) Warn(msg string, information Information) {
	l.logger.Warn().Str("location", getLocation(2)).Fields(information).Msg(msg)
}

func (l *ZeroLog) Error(msg string, err error, information Information) {
	l.logger.Err(err).Str("location", getLocation(2)).Fields(information).Msg(msg)
}

// getLocation returns a string that contains the file and
// the line in the file where the error occurred
func getLocation(calls int) (location string) {
	_, file, line, ok := runtime.Caller(calls)
	path := strings.Split(file, "/")
	if ok {
		location = fmt.Sprintf("%s:%d", path[len(path)-1], line)
	}
	return
}
