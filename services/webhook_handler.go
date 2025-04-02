package services

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github_webhook/domain"
	"github_webhook/utils"

	"github.com/gin-gonic/gin"
)

type WebhookHandler struct {
	discordService *utils.DiscordService
}

func NewWebhookHandler(discordService *utils.DiscordService) *WebhookHandler {
	return &WebhookHandler{
		discordService: discordService,
	}
}

func HandleGitHubWebhook(c *gin.Context) {
	handler := NewWebhookHandler(utils.NewDiscordService())

	eventType := c.GetHeader("X-GitHub-Event")
	deliveryID := c.GetHeader("X-GitHub-Delivery")

	log.Printf("Evento recibido: %s, ID: %s", eventType, deliveryID)

	payload, err := c.GetRawData()
	if err != nil {
		log.Printf("Error al leer el payload: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Error al leer el payload"})
		return
	}

	switch eventType {
	case "ping":
		handlePing(c)
	case "pull_request":
		handlePullRequest(c, payload, handler)
	case "workflow_run":
		handleWorkflowRun(c, payload, handler)
	default:
		log.Printf("Evento no manejado: %s", eventType)
		c.JSON(http.StatusOK, gin.H{"message": "Evento no manejado"})
	}
}

func handlePing(c *gin.Context) {
	log.Printf("Ping recibido correctamente")
	c.JSON(http.StatusOK, gin.H{"message": "Ping recibido correctamente"})
}

func handlePullRequest(c *gin.Context, payload []byte, handler *WebhookHandler) {
	var prEvent domain.PullRequestEvent
	if err := json.Unmarshal(payload, &prEvent); err != nil {
		log.Printf("Error al deserializar el evento: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Error al procesar el evento"})
		return
	}

	log.Printf("Procesando Pull Request: %s, Acción: %s", prEvent.PullRequest.Title, prEvent.Action)

	// Solo procesar eventos específicos
	if !isRelevantPRAction(prEvent.Action) {
		log.Printf("Acción de PR no relevante: %s", prEvent.Action)
		c.JSON(http.StatusOK, gin.H{"message": "Acción de PR no relevante"})
		return
	}

	message := formatPullRequestMessage(prEvent)
	log.Printf("Enviando mensaje a Discord (Dev): %s", message)

	if err := handler.discordService.SendDevMessage(message); err != nil {
		log.Printf("Error al enviar mensaje a Discord: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al enviar notificación"})
		return
	}

	log.Printf("Mensaje enviado exitosamente a Discord")
	c.JSON(http.StatusOK, gin.H{"message": "Evento procesado correctamente"})
}

func handleWorkflowRun(c *gin.Context, payload []byte, handler *WebhookHandler) {
	var workflowEvent domain.WorkflowRunEvent
	if err := json.Unmarshal(payload, &workflowEvent); err != nil {
		log.Printf("Error al deserializar el evento de workflow: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Error al procesar el evento"})
		return
	}

	log.Printf("Procesando Workflow: %s, Estado: %s", workflowEvent.Workflow.Name, workflowEvent.Status)

	message := formatWorkflowMessage(workflowEvent)
	log.Printf("Enviando mensaje a Discord (Test): %s", message)

	if err := handler.discordService.SendTestMessage(message); err != nil {
		log.Printf("Error al enviar mensaje a Discord: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al enviar notificación"})
		return
	}

	log.Printf("Mensaje enviado exitosamente a Discord")
	c.JSON(http.StatusOK, gin.H{"message": "Evento de workflow procesado correctamente"})
}

func isRelevantPRAction(action string) bool {
	relevantActions := map[string]bool{
		"opened":           true,
		"reopened":         true,
		"ready_for_review": true,
		"closed":           true,
	}
	return relevantActions[action]
}

func formatPullRequestMessage(event domain.PullRequestEvent) string {
	action := event.Action
	pr := event.PullRequest
	repo := event.Repository

	status := "❌"
	if action == "closed" && pr.Merged {
		status = "✅"
		action = "merged"
	}

	return fmt.Sprintf(
		"**%s Pull Request** %s\n"+
			"Repositorio: %s\n"+
			"Acción: %s\n"+
			"Título: %s\n"+
			"URL: %s\n"+
			"Estado: %s",
		status,
		pr.Title,
		repo.FullName,
		action,
		pr.Title,
		pr.HTMLURL,
		pr.State,
	)
}

func formatWorkflowMessage(event domain.WorkflowRunEvent) string {
	status := "❌"
	if event.Conclusion == "success" {
		status = "✅"
	}

	return fmt.Sprintf(
		"**Workflow Run** %s\n"+
			"Repositorio: %s\n"+
			"Workflow: %s\n"+
			"Estado: %s\n"+
			"Conclusión: %s\n"+
			"URL: %s",
		status,
		event.Repository.FullName,
		event.Workflow.Name,
		event.Status,
		event.Conclusion,
		event.HTMLURL,
	)
}
