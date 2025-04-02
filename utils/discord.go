package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
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

	log.Printf("Configuración de Discord:")
	log.Printf("DISCORD_WEBHOOK_DESARROLO: %s", devWebhook)
	log.Printf("DISCORD_WEBHOOK_PRUEBAS: %s", testWebhook)

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
	log.Printf("Intentando enviar mensaje al canal de desarrollo: %s", content)
	return s.sendMessage(s.devWebhookURL, content)
}

func (s *DiscordService) SendTestMessage(content string) error {
	if s.testWebhookURL == "" {
		return fmt.Errorf("DISCORD_WEBHOOK_PRUEBAS no está configurado")
	}
	log.Printf("Intentando enviar mensaje al canal de pruebas: %s", content)
	return s.sendMessage(s.testWebhookURL, content)
}

func (s *DiscordService) sendMessage(webhookURL, content string) error {
	message := DiscordMessage{
		Content: content,
	}

	jsonData, err := json.Marshal(message)
	if err != nil {
		log.Printf("Error al serializar el mensaje: %v", err)
		return fmt.Errorf("error al serializar el mensaje: %v", err)
	}

	log.Printf("Enviando mensaje a Discord:")
	log.Printf("- URL: %s", webhookURL)
	log.Printf("- Contenido: %s", string(jsonData))

	req, err := http.NewRequest("POST", webhookURL, bytes.NewBuffer(jsonData))
	if err != nil {
		log.Printf("Error al crear la petición: %v", err)
		return fmt.Errorf("error al crear la petición: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Error al enviar la petición: %v", err)
		return fmt.Errorf("error al enviar mensaje a Discord: %v", err)
	}
	defer resp.Body.Close()

	log.Printf("Respuesta de Discord:")
	log.Printf("- Status Code: %d", resp.StatusCode)
	log.Printf("- Status: %s", resp.Status)

	if resp.StatusCode != http.StatusNoContent && resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		log.Printf("Cuerpo de la respuesta de error: %s", string(body))
		return fmt.Errorf("error en la respuesta de Discord: %d - %s", resp.StatusCode, resp.Status)
	}

	log.Printf("Mensaje enviado exitosamente a Discord")
	return nil
}
