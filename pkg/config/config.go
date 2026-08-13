package config

import (
	"os"
)

const (
	ConfigAppEnv     = "APP_ENV"
	ConfigLogLevel   = "LOG_LEVEL"
	ConfigServerMode = "SERVER_MODE"

	ConfigAppHost     = "APP_HOST"
	ConfigAppScheme   = "APP_SCHEME"
	ConfigAppBasePath = "APP_BASE_PATH"
	ConfigAppName     = "APP_NAME"
	ConfigAppPort     = "APP_PORT"

	ConfigDBUser     = "DATABASE_USER"
	ConfigDBPassword = "DATABASE_PASSWORD"
	ConfigDBHost     = "DATABASE_HOST"
	ConfigDBPort     = "DATABASE_PORT"
	ConfigDBName     = "DATABASE_NAME"

	ConfigRMQVirtualHost = "RABBIT_VIRT_HOST"
	ConfigRMQHost        = "RABBIT_HOST"
	ConfigRMQPort        = "RABBIT_PORT"
	ConfigRMQUser        = "RABBIT_USER"
	ConfigRMQPassword    = "RABBIT_PASSWORD"
)

func GetAll() map[string]string {
	configs := make(map[string]string)

	configs[ConfigAppEnv] = getFromEnvOrDefault(ConfigAppEnv, "prod")
	configs[ConfigLogLevel] = getFromEnvOrDefault(ConfigLogLevel, "1")
	configs[ConfigServerMode] = getFromEnvOrDefault(ConfigServerMode, "release")

	configs[ConfigAppHost] = getFromEnvOrDefault(ConfigAppHost, "http://localhost")
	configs[ConfigAppScheme] = getFromEnvOrDefault(ConfigAppScheme, "http")
	configs[ConfigAppBasePath] = getFromEnvOrDefault(ConfigAppBasePath, "/api")
	configs[ConfigAppName] = getFromEnvOrDefault(ConfigAppName, "scootgo")
	configs[ConfigAppPort] = getFromEnvOrDefault(ConfigAppPort, "8080")

	configs[ConfigDBUser] = getFromEnvOrDefault(ConfigDBUser, "root")
	configs[ConfigDBPassword] = getFromEnvOrDefault(ConfigDBPassword, "toor")
	configs[ConfigDBHost] = getFromEnvOrDefault(ConfigDBHost, "localhost")
	configs[ConfigDBPort] = getFromEnvOrDefault(ConfigDBPort, "3306")
	configs[ConfigDBName] = getFromEnvOrDefault(ConfigDBName, "local")

	configs[ConfigRMQVirtualHost] = getFromEnvOrDefault(ConfigRMQVirtualHost, "local_virtual")
	configs[ConfigRMQHost] = getFromEnvOrDefault(ConfigRMQHost, "localhost")
	configs[ConfigRMQPort] = getFromEnvOrDefault(ConfigRMQPort, "5672")
	configs[ConfigRMQUser] = getFromEnvOrDefault(ConfigRMQUser, "guest")
	configs[ConfigRMQPassword] = getFromEnvOrDefault(ConfigRMQPassword, "guest")

	return configs
}

func getFromEnvOrDefault(key, defaultVal string) string {
	value, found := os.LookupEnv(key)
	if !found || value == "" {
		return defaultVal
	}

	return value
}
