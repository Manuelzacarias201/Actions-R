package services

import (
	"bytes"
	"fmt"
	"io"
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

func (h *WebhookHandler) HandleGitHubWebhook(c *gin.Context) {
	eventType := c.GetHeader("X-GitHub-Event")
	deliveryID := c.GetHeader("X-GitHub-Delivery")
	log.Printf("Evento recibido: %s, Delivery ID: %s", eventType, deliveryID)

	// Log headers
	log.Printf("Headers recibidos: %v", c.Request.Header)

	// Log payload
	body, _ := io.ReadAll(c.Request.Body)
	log.Printf("Payload recibido: %s", string(body))
	c.Request.Body = io.NopCloser(bytes.NewBuffer(body))

	var event interface{}
	switch eventType {
	case "pull_request":
		event = &domain.PullRequestEvent{}
	case "workflow_run":
		event = &domain.WorkflowRunEvent{}
	case "workflow_job":
		event = &domain.WorkflowJobEvent{}
	default:
		log.Printf("Evento no manejado: %s", eventType)
		c.JSON(http.StatusOK, gin.H{"message": "Evento no manejado"})
		return
	}

	if err := c.ShouldBindJSON(event); err != nil {
		log.Printf("Error al decodificar el evento: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	switch e := event.(type) {
	case *domain.PullRequestEvent:
		h.handlePullRequestEvent(e)
	case *domain.WorkflowRunEvent:
		h.handleWorkflowRunEvent(e)
	case *domain.WorkflowJobEvent:
		h.handleWorkflowJobEvent(e)
	}

	c.JSON(http.StatusOK, gin.H{"message": "Evento procesado"})
}

func handlePing(c *gin.Context) {
	log.Printf("Ping recibido correctamente")
	c.JSON(http.StatusOK, gin.H{"message": "Ping recibido correctamente"})
}

func (h *WebhookHandler) handlePullRequestEvent(event *domain.PullRequestEvent) {
	log.Printf("Procesando Pull Request: %s, Acción: %s", event.PullRequest.Title, event.Action)

	// Solo procesar eventos específicos
	if !isRelevantPRAction(event.Action) {
		log.Printf("Acción de PR no relevante: %s", event.Action)
		return
	}

	message := formatPullRequestMessage(*event)
	log.Printf("Enviando mensaje a Discord (Dev): %s", message)

	if err := h.discordService.SendDevMessage(message); err != nil {
		log.Printf("Error al enviar mensaje a Discord: %v", err)
	}
}

func (h *WebhookHandler) handleWorkflowRunEvent(event *domain.WorkflowRunEvent) {
	log.Printf("Detalles del evento de workflow:")
	log.Printf("- Nombre del workflow: %s", event.Workflow.Name)
	log.Printf("- Estado: %s", event.Status)
	log.Printf("- Conclusión: %s", event.Conclusion)
	log.Printf("- Repositorio: %s", event.Repository.FullName)
	log.Printf("- URL: %s", event.HTMLURL)
	log.Printf("- Evento que lo activó: %s", event.Event)
	log.Printf("- Rama: %s", event.HeadBranch)

	// Solo enviar notificación cuando el workflow termine
	if event.Status != "completed" {
		log.Printf("Workflow aún no ha terminado, estado actual: %s", event.Status)
		return
	}

	message := formatWorkflowMessage(*event)
	log.Printf("Enviando mensaje a Discord (Test): %s", message)

	if err := h.discordService.SendTestMessage(message); err != nil {
		log.Printf("Error al enviar mensaje a Discord: %v", err)
	}
}

func (h *WebhookHandler) handleWorkflowJobEvent(event *domain.WorkflowJobEvent) {
	log.Printf("Procesando evento workflow_job: %s", event.Action)

	// Solo enviamos notificaciones cuando el job se completa
	if event.Action == "completed" {
		emoji := "✅"
		if event.WorkflowJob.Conclusion != "success" {
			emoji = "❌"
		}

		message := fmt.Sprintf("%s **Workflow Job Completado**\n\n"+
			"**Repositorio:** %s\n"+
			"**Job:** %s\n"+
			"**Estado:** %s\n"+
			"**Conclusión:** %s\n"+
			"**URL:** %s\n"+
			"**Iniciado:** %s\n"+
			"**Completado:** %s",
			emoji,
			event.Repository.FullName,
			event.WorkflowJob.Name,
			event.WorkflowJob.Status,
			event.WorkflowJob.Conclusion,
			event.WorkflowJob.HTMLURL,
			event.WorkflowJob.StartedAt,
			event.WorkflowJob.CompletedAt)

		if err := h.discordService.SendTestMessage(message); err != nil {
			log.Printf("Error al enviar mensaje a Discord: %v", err)
		}
	}
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

	emoji := "❌"
	if action == "closed" && pr.Merged {
		emoji = "✅"
		action = "merged"
	}

	return fmt.Sprintf(
		"**%s Pull Request** %s\n"+
			"Repositorio: %s\n"+
			"Acción: %s\n"+
			"Título: %s\n"+
			"URL: %s\n"+
			"Estado: %s",
		emoji,
		pr.Title,
		repo.FullName,
		action,
		pr.Title,
		pr.HTMLURL,
		pr.State,
	)
}

func formatWorkflowMessage(event domain.WorkflowRunEvent) string {
	emoji := "✅"
	if event.Conclusion == "failure" {
		emoji = "❌"
	} else if event.Conclusion == "success" {
		emoji = "✅"
	} else {
		emoji = "⚠️"
	}

	status := "✅ Éxito"
	if event.Conclusion == "failure" {
		status = "❌ Fallo"
	} else if event.Conclusion == "success" {
		status = "✅ Éxito"
	} else {
		status = "⚠️ Otro"
	}

	return fmt.Sprintf(
		"**%s Workflow Run**\n"+
			"Repositorio: %s\n"+
			"Workflow: %s\n"+
			"Estado: %s\n"+
			"Conclusión: %s\n"+
			"URL: %s",
		emoji,
		event.Repository.FullName,
		event.Workflow.Name,
		status,
		event.Conclusion,
		event.HTMLURL,
	)
}
