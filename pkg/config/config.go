package config

import (
	"forum/pkg/serror"
	"time"
)

// default constants
const (
	defaultHTTPServerHost        = "localhost"
	defaultHTTPServerPort        = "8080"
	defaultHTTPServeIdleTime     = 60 * time.Second
	defaultHTTPServerWriteTime   = 30 * time.Second
	defaultHTTPServerReadTime    = 30 * time.Second
	defaultHTTPServerMaxHeaderMb = 20 << 20 // 20mb

	defaultLoggerLevel     = -4 // Debug Level
	defaultLoggerSourceKey = true
	defaultLoggerOutput    = "stdout"
	defaultLoggerHandler   = "json"
	// nanosecond = 1
	// microsecond = 1000 nano
	// milisecond= 1000 micro
	// second = 1000
	// minute = 60
)

// Config structure
type (
	HTTPServer struct {
		Host         string
		Port         string
		IdleTimeout  time.Duration
		WriteTimeout time.Duration
		ReadTimeout  time.Duration
		MaxHeaderMb  int
	}

	Logger struct {
		Level     int
		SourceKey bool
		Output    string
		Handler   string
	}

	Config struct {
		ServiceName string
		HTTPServer
		Logger
	}
)

// InitConfig ...
func InitConfig(cfg *Config) (*Config, error) {
	// setup

	// someVariable := "some"

	// process := []byte{}

	// env - process

	// user space

	//system - root -$PATH - exec : .env
	cfg, err := setup(cfg)
	if err != nil {
		return nil, err
	}

	return cfg, nil

}

// Переменная окружения — это переменная, значение которой используется операционной системой или приложениями для конфигурации работы программного обеспечения в конкретной среде.
func setup(cfg *Config) (cnfg *Config, err error) {
	//setup method <
	if cfg == nil {
		return nil, serror.ErrConfigStructureIsNill
	}

	// parsing
	// .env ^^^, json, config constants, -> validatation
	// set conf paras
	// or populate defaults

}

// config, logger,model -> entity, repository. connection to database, query

// || server, router, meddleware, auth
