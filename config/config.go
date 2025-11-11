package config

const (
	// DefaultPort for the TCP server
	DefaultPort = "8080"
)

// Config holds application configuration
type Config struct {
	// Server
	Port      string
	JWTSecret string

	// Logging
	LogLevel      string
	LogDir        string
	LogMaxSize    int
	LogMaxBackups int
	LogMaxAge     int
	LogCompress   bool
	LogConsole    bool
}

// Load loads configuration from .env file and environment variables
func Load() *Config {
	// Load .env file
	LoadEnv()

	cfg := &Config{
		// Server
		Port:      GetEnv("PORT", DefaultPort),
		JWTSecret: MustGetEnv("JWT_SECRET"),

		// Logging
		LogLevel:      GetEnv("LOG_LEVEL", "info"),
		LogDir:        GetEnv("LOG_DIR", "./logs"),
		LogMaxSize:    GetEnvInt("LOG_MAX_SIZE", 100),
		LogMaxBackups: GetEnvInt("LOG_MAX_BACKUPS", 5),
		LogMaxAge:     GetEnvInt("LOG_MAX_AGE", 30),
		LogCompress:   GetEnvBool("LOG_COMPRESS", true),
		LogConsole:    GetEnvBool("LOG_CONSOLE", true),
	}

	return cfg
}
