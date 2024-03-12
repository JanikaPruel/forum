package app

import (
	"fmt"
	"forum/pkg/config"
)

func Run(cfg *config.Config) {

	// Init router

	fmt.Printf("\n\n Server runnig on -> http://%s:%s\n\n Enter Ctrl + C for stop application", cfg.HTTPServer.Host, cfg.HTTPServer.Port)
	// Perfect end our application
}
