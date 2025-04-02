package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

type DiscordService struct {
	devWebhookURL  string
	testWebhookURL string
}

type DiscordMessage struct {
	Content string `json:"content"`
}

func NewDiscordService() *DiscordService {
	return &DiscordService{
		devWebhookURL:  os.Getenv("DISCORD_DEV_WEBHOOK_URL"),
		testWebhookURL: os.Getenv("DISCORD_TEST_WEBHOOK_URL"),
	}
}

func (s *DiscordService) SendDevMessage(content string) error {
	return s.sendMessage(s.devWebhookURL, content)
}

func (s *DiscordService) SendTestMessage(content string) error {
	return s.sendMessage(s.testWebhookURL, content)
}

func (s *DiscordService) sendMessage(webhookURL, content string) error {
	message := DiscordMessage{
		Content: content,
	}

	jsonData, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("error al serializar el mensaje: %v", err)
	}

	resp, err := http.Post(
		webhookURL,
		"application/json",
		bytes.NewBuffer(jsonData),
	)
	if err != nil {
		return fmt.Errorf("error al enviar mensaje a Discord: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent && resp.StatusCode != http.StatusOK {
		return fmt.Errorf("error en la respuesta de Discord: %d", resp.StatusCode)
	}

	return nil
}
