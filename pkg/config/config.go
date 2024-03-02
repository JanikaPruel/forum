package config

import "time"

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
func InitConfig(cfg *Config, err error) {
	// setup
	//setup method <
	// parsing
	// json, validatation
	// set conf paras
	// or populate defaults

}
