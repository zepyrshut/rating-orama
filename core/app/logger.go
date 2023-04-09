package app

import (
	"golang.org/x/exp/slog"
	"io"
	"os"
	"time"
)

func newStructuredLogger() *slog.Logger {
	logFile, _ := os.OpenFile(generateLogFileName(), os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	multiWriter := io.MultiWriter(os.Stdout, logFile)

	return slog.New(slog.NewJSONHandler(multiWriter))
}

func generateLogFileName() string {
	currentTime := time.Now()
	return "logs/" + currentTime.Format("2006-01-02") + ".log"
}
