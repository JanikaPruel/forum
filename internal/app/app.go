package app

import (
	"context"
	"fmt"
	"forum/internal/controller/router"
	"forum/internal/controller/server"
	"forum/pkg/config"
	"forum/pkg/sqlite"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
)

func Run(cfg *config.Config) error {

	// Init router
	router := router.New()
	router.InitRouter()

	// Perfect end our application

	// prepare database for connection
	// connection to database
	db, err := sqlite.ConnectionToDB(cfg)
	if err != nil {
		slog.Error(err.Error())
		return err
	}

	// db.ConnectionToDB(cfg)
	if err := db.SQLite.Ping(); err != nil {
		slog.Error(err.Error())
		return err
	}

	slog.Info("Susseccfull to connection data base")

	srv := server.New(cfg, router)

	slog.Debug("http server successfully created")

	//running the http server + -> for{ infinity logic}
	fmt.Printf("\n\nServer running on -> http://%s:%s\n\nPress Ctrl+C for stop\n\n", cfg.HTTPServer.Host, cfg.HTTPServer.Port)
	go func() {
		if err := srv.Run(); err != nil {
			slog.Error(err.Error())
		}
	}()

	// getting API -> webapi <-controller ---

	// gracefull shutdown -> stopping the http server +
	sigChan := make(chan os.Signal, 1)

	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGTERM)

	// logging INFO
	sig := <-sigChan
	slog.Info("Some crazy developer termination our app, with: ", "syscall", sig.String())

	if err := srv.Shutdown(context.Background()); err != nil {
		slog.Error(err.Error())

		return err
	}

	slog.Info("Server successfull stopped")
	if err := db.SQLite.Close(); err != nil {
		slog.Error(err.Error())
		return err
	}
	slog.Info("Database connection successfull closed")
	return nil
	// middleware

	// controller

	// http server

	//start http server

	// syscall signal
	// close database connnection
	// stop http server

}
