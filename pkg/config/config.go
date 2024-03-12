package config

import (
	"encoding/json"
	"errors"
	"forum/pkg/serror"
	"os"
	"strings"
	"time"
)

// default constants
const (
	defaultServiceName = "forum"

	defaultHTTPServerHost         = "localhost"
	defaultHTTPServerPort         = "8080"
	defaultHTTPServeIdleTimeout   = 60 * time.Second
	defaultHTTPServerWriteTimeout = 30 * time.Second
	defaultHTTPServerReadTimeout  = 30 * time.Second
	defaultHTTPServerMaxHeaderMb  = 20 << 20 // 20mb

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
		Host         string        `json:"host" env:"SERVER_HOST"`
		Port         string        `json:"port"  env:"SERVER_PORT"`
		IdleTimeout  time.Duration `json:"idle_time"`
		WriteTimeout time.Duration `json:"write_time"`
		ReadTimeout  time.Duration `json:"read_time"`
		MaxHeaderMb  int           `json:"max_header_mb"`
	}

	LoggerConfig struct {
		Level     int    `json:"level"`
		SourceKey bool   `json:"source_key"`
		Output    string `json:"output"`
		Handler   string `json:"handler"`
	}

	Config struct {
		ServiceName  string `json:"service_name"`
		HTTPServer   `json:"http_server"`
		LoggerConfig `json:"logger"`
	}
)

// Config path
const (
	configDir  = "configs"
	configFile = "default.json"
)

// InitConfig ...
func InitConfig() (*Config, error) {
	cfg := &Config{}

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
func setup(cfg *Config) (*Config, error) {
	//setup method <
	if cfg == nil {
		return nil, serror.ErrConfigStructureIsNill
	}
	configFilePath := strings.Join([]string{configDir, configFile}, "/")

	populateDefaults(cfg)

	cfg, err := parseConfigFileAndSetConfigParams(configFilePath, cfg)
	if err != nil {
		return nil, err
	}

	// setConfigParamsFromEnv
	cfg, err = setConfigParamsFromEnv(cfg)
	if err != nil {
		return nil, err
	}

	// parsing
	// .env ^^^, json, config constants, -> validatation
	// set conf paras
	// or populate defaults

	return cfg, nil
}

// Чувствительный данные их храним в переменных окружения ОС.
// setConfigParamsFromEnv
func setConfigParamsFromEnv(cfg *Config) (*Config, error) {
	_, exist := os.LookupEnv("SERVER_HOST")
	if !exist {
		return nil, errors.New("error, env not exists")
	}

	cfg.HTTPServer.Host = os.Getenv("SERVER_HOST")
	cfg.HTTPServer.Port = os.Getenv("SERVER_PORT")

	return cfg, nil
}

// parseConfigFileAndSetConfigParams
func parseConfigFileAndSetConfigParams(filePath string, cfg *Config) (*Config, error) {
	if filePath == "" {
		return nil, errors.New("error, invalid config file path")
	}

	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(data, cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}

// populateDefaults
func populateDefaults(cfg *Config) {
	// service
	cfg.ServiceName = defaultServiceName
	// HTTPServer
	cfg.HTTPServer.Host = defaultHTTPServerHost
	cfg.HTTPServer.Port = defaultHTTPServerPort
	cfg.HTTPServer.IdleTimeout = defaultHTTPServeIdleTimeout
	cfg.HTTPServer.WriteTimeout = defaultHTTPServerWriteTimeout
	cfg.HTTPServer.ReadTimeout = defaultHTTPServerReadTimeout
	// Logger
	cfg.LoggerConfig.Level = defaultLoggerLevel
	cfg.LoggerConfig.SourceKey = defaultLoggerSourceKey
	cfg.LoggerConfig.Output = defaultLoggerOutput
	cfg.LoggerConfig.Handler = defaultLoggerHandler
}

// config, logger,model -> entity, repository. connection to database, query

// || server, router, meddleware, auth
