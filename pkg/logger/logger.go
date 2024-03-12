package logger

import (
	"errors"
	"forum/pkg/config"
	"log/slog"
	"os"
)

// InitLogger return a new customized slog logger instance
func InitLogger(cfg *config.Config) (logger *slog.Logger, err error) {
	// validate config params
	if !isValidLoggerConfigParams(cfg) {
		return nil, errors.New("error, invalid config params")
	}

	// choice output = stdout or
	// choice handler

	logger = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level:     slog.Level(cfg.LoggerConfig.Level),
		AddSource: cfg.LoggerConfig.SourceKey,
	}).WithAttrs([]slog.Attr{slog.String("service_name", cfg.ServiceName)}))

	slog.SetDefault(logger)

	return logger, nil
}

// isValidLoggerConfigParams
func isValidLoggerConfigParams(cfg *config.Config) (valid bool) {
	if cfg == nil {
		return false
	}

	if !(cfg.LoggerConfig.Handler == "json" || cfg.LoggerConfig.Handler == "text") {
		return false
	} else if !(cfg.LoggerConfig.Output == "output" || cfg.LoggerConfig.Output == "file") {
		return false
	}
	return true
}
