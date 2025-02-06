package logger

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/lucasmendesl/beerus/version"
)

func Create(logFormat, logLevel string) (*slog.Logger, error) {
	handler, err := createHandler(logFormat)
	if err != nil {
		return nil, err
	}

	level := new(slog.LevelVar)
	if err := level.UnmarshalText([]byte(logLevel)); err != nil {
		return nil, fmt.Errorf("invalid log level: %w", err)
	}

	slog.SetLogLoggerLevel(level.Level())
	return slog.New(handler).With("service", "beerus", "version", version.Version), nil
}

func createHandler(logFormatter string) (slog.Handler, error) {
	switch logFormatter {
	case "json":
		return slog.NewJSONHandler(os.Stdout, nil), nil
	case "text":
		return slog.NewTextHandler(os.Stdout, nil), nil
	default:
		return nil, fmt.Errorf("invalid log formatter: %s", logFormatter)
	}
}
