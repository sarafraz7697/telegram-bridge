package main

import (
	"net/http"

	"telegram-bridge/config"
	"telegram-bridge/handlers"
	"telegram-bridge/logger"
	"telegram-bridge/router"
	"telegram-bridge/telegram"
)

func main() {
	// Load configuration
	cfg := config.Load()

	// Initialize logger
	logCfg := &logger.Config{
		Level:      cfg.LogLevel,
		LogDir:     cfg.LogDir,
		MaxSize:    cfg.LogMaxSize,
		MaxBackups: cfg.LogMaxBackups,
		MaxAge:     cfg.LogMaxAge,
		Compress:   cfg.LogCompress,
		Console:    cfg.LogConsole,
	}

	log, err := logger.New("telegram-bridge", logCfg)
	if err != nil {
		panic("Failed to initialize logger: " + err.Error())
	}
	defer log.Sync()

	// Initialize dependencies
	telegramClient := telegram.NewClient(log)
	handler := handlers.NewHandler(telegramClient, log)

	// Setup routes
	mux := router.SetupRoutes(handler, cfg.AuthToken, log)

	// Log startup information
	log.Info("Starting Telegram API Bridge",
		logger.String("port", cfg.Port),
		logger.String("log_level", cfg.LogLevel))
	log.Info("Available endpoints",
		logger.Strings("endpoints", []string{
			"POST /send - Send message to Telegram",
			"GET /health - Health check",
		}))

	// Start server
	log.Info("Server listening", logger.String("port", cfg.Port))
	if err := http.ListenAndServe(":"+cfg.Port, mux); err != nil {
		log.Fatal("Server failed to start", logger.ZError(err))
	}
}
