package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
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
	devWebhook := os.Getenv("DISCORD_WEBHOOK_DESARROLO")
	testWebhook := os.Getenv("DISCORD_WEBHOOK_PRUEBAS")

	if devWebhook == "" {
		log.Printf("ADVERTENCIA: DISCORD_WEBHOOK_DESARROLO no está configurado")
	}
	if testWebhook == "" {
		log.Printf("ADVERTENCIA: DISCORD_WEBHOOK_PRUEBAS no está configurado")
	}

	return &DiscordService{
		devWebhookURL:  devWebhook,
		testWebhookURL: testWebhook,
	}
}

func (s *DiscordService) SendDevMessage(content string) error {
	if s.devWebhookURL == "" {
		return fmt.Errorf("DISCORD_WEBHOOK_DESARROLO no está configurado")
	}
	return s.sendMessage(s.devWebhookURL, content)
}

func (s *DiscordService) SendTestMessage(content string) error {
	if s.testWebhookURL == "" {
		return fmt.Errorf("DISCORD_WEBHOOK_PRUEBAS no está configurado")
	}
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

	log.Printf("Enviando mensaje a Discord URL: %s", webhookURL)

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

	log.Printf("Mensaje enviado exitosamente a Discord")
	return nil
}
