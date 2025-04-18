package config

import (
	"os"
)

type Config struct {
	port               string
	discordDevWebhook  string
	discordTestWebhook string
}
//oo
func NewConfig() *Config {
	return &Config{
		port:               getEnvOrDefault("PORT", "8080"),
		discordDevWebhook:  getEnvOrDefault("DISCORD_WEBHOOK_DESARROLO", ""),
		discordTestWebhook: getEnvOrDefault("DISCORD_WEBHOOK_PRUEBAS", ""),
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
