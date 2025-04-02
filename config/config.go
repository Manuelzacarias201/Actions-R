package config

import (
	"os"
)

type Config struct {
	port               string
	discordDevWebhook  string
	discordTestWebhook string
}

func NewConfig() *Config {
	return &Config{
		port:               getEnvOrDefault("PORT", "8080"),
		discordDevWebhook:  getEnvOrDefault("DISCORD_DEV_WEBHOOK_URL", ""),
		discordTestWebhook: getEnvOrDefault("DISCORD_TEST_WEBHOOK_URL", ""),
	}
}

func (c *Config) GetPort() string {
	return c.port
}

func (c *Config) GetDiscordDevWebhook() string {
	return c.discordDevWebhook
}

func (c *Config) GetDiscordTestWebhook() string {
	return c.discordTestWebhook
}

func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
