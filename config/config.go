package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	TelegramBotToken string
	TelegramChatID   string
	MySQLHost        string
	MySQLPort        string
	MySQLUser        string
	MySQLPassword    string
	MySQLDatabase    string
	LogLevel         string
}

func Load() *Config {
	err := godotenv.Load()
	if err != nil {
		log.Println("Warning: .env file not found, using environment variables")
	}

	return &Config{
		TelegramBotToken: getEnv("TELEGRAM_BOT_TOKEN", ""),
		TelegramChatID:   getEnv("TELEGRAM_CHAT_ID", ""),
		MySQLHost:        getEnv("MYSQL_HOST", "localhost"),
		MySQLPort:        getEnv("MYSQL_PORT", "3306"),
		MySQLUser:        getEnv("MYSQL_USER", "root"),
		MySQLPassword:    getEnv("MYSQL_PASSWORD", ""),
		MySQLDatabase:    getEnv("MYSQL_DATABASE", "exchange_symbols"),
		LogLevel:         getEnv("LOG_LEVEL", "info"),
	}
}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
