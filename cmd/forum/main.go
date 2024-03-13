package main

import (
	"forum/internal/app"
	"forum/pkg/config"
	"forum/pkg/logger"
	"log/slog"
)

// Если входящие данные критичны то мы ОБЯЗАНЫ проверять их!!! RED CASE - VALIDATE DATE. Могут прийти не корректные данные и наша
// Программа может отработать не корректно или внедрить в наш сервис скрипт и сложить наш прод. Атака называется SqlScript
func main() {
	// InitConfig
	cfg, err := config.InitConfig()
	if err != nil {
		slog.Error("Ny vse priehali davai zanovo init config")
	}

	// Logger
	_, err = logger.InitLogger(cfg) // cfg = nil
	if err != nil {
		slog.Error("Ny vse priehali davai zanovo logger config")
	}

	// Debug log for test
	slog.Debug("OOOOO vse harasho, normalno edem. Viy viy viy")

	// app -> Run()
	if err := app.Run(cfg); err != nil {
		slog.Debug("OOOO syeta tormozi davai check app run")
	}
}

// sqlite3
