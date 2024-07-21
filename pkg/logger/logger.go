package logger

import (
	"fmt"
	"log/slog"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

type Logger struct {
	logger *slog.Logger
}

func New(debug bool) *Logger {
	level := slog.LevelInfo
	if debug {
		level = slog.LevelDebug
	}

	logFile, err := tea.LogToFile("debug.log", "")
	if err != nil {
		fmt.Printf("failed to create log file: %v\n", err)
		os.Exit(1)
	}

	handler := slog.NewJSONHandler(logFile, &slog.HandlerOptions{
		Level: level,
	})

	l := &Logger{
		logger: slog.New(handler),
	}

	return l
}

func (l *Logger) Debug(msg string, args ...interface{}) {
	l.logger.Debug(msg, args...)
}

func (l *Logger) Info(msg string, args ...interface{}) {
	l.logger.Info(msg, args...)
}

func (l *Logger) Warn(msg string, args ...interface{}) {
	l.logger.Warn(msg, args...)
}

func (l *Logger) Error(msg string, args ...interface{}) {
	l.logger.Error(msg, args...)
}
