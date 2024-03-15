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

	// TODO: bigfix

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

	if cfg.LoggerConfig.Level < int(slog.LevelDebug) || cfg.LoggerConfig.Level > int(slog.LevelError) {
		return false
	}

	if !(cfg.LoggerConfig.Handler == "json" || cfg.LoggerConfig.Handler == "text") {
		return false
	}

	if !(cfg.LoggerConfig.Output == "stdout" || cfg.LoggerConfig.Output == "file") { // bug
		return false

	}
	return true
}
