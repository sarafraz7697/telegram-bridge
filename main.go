package main

import (
	"log"
	"os"

	"telegram-bridge/config"
	"telegram-bridge/logger"
	"telegram-bridge/server"
)

func main() {
	// Load configuration
	cfg := config.Load()

	// Initialize logger
	logInstance, err := logger.New("telegram-bridge", &logger.Config{
		Level:      cfg.LogLevel,
		LogDir:     cfg.LogDir,
		MaxSize:    cfg.LogMaxSize,
		MaxBackups: cfg.LogMaxBackups,
		MaxAge:     cfg.LogMaxAge,
		Compress:   cfg.LogCompress,
		Console:    cfg.LogConsole,
	})
	if err != nil {
		log.Fatalf("Failed to initialize logger: %v", err)
	}
	defer logInstance.Sync()

	// Set global logger
	logger.SetGlobal(logInstance)

	logInstance.Info("Starting Telegram Bridge Server")

	// Initialize and start TCP server
	tcpServer := server.NewTCPServer(cfg, logInstance)

	logInstance.Info("Server configuration",
		logger.String("port", cfg.Port),
		logger.String("log_level", cfg.LogLevel),
		logger.String("log_dir", cfg.LogDir))

	// Start the server (blocking call)
	if err := tcpServer.Start(); err != nil {
		logInstance.Fatal("Failed to start TCP server", logger.String("error", err.Error()))
		os.Exit(1)
	}
}
