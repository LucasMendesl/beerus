package logger

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/lucasmendesl/beerus/config"
	"github.com/lucasmendesl/beerus/version"
)

func Create(config config.Logging) (*slog.Logger, error) {
	level := new(slog.LevelVar)
	if err := level.UnmarshalText([]byte(config.Level)); err != nil {
		return nil, fmt.Errorf("invalid log level: %w", err)
	}

	handler, err := createHandler(config.Format, level)
	if err != nil {
		return nil, err
	}

	return slog.New(handler).With("service", "beerus", "version", version.Version), nil
}

func createHandler(logFormatter string, level *slog.LevelVar) (slog.Handler, error) {
	handlerOpts := &slog.HandlerOptions{
		Level: level,
	}

	switch logFormatter {
	case "json":
		return slog.NewJSONHandler(os.Stdout, handlerOpts), nil
	case "text":
		return slog.NewTextHandler(os.Stdout, handlerOpts), nil
	default:
		return nil, fmt.Errorf("invalid log formatter: %s", logFormatter)
	}
}
