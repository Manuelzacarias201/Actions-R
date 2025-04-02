package http

import (
	"github.com/actions/internal/core/domain"
	"github.com/actions/internal/core/ports"
	"github.com/gin-gonic/gin"
)

type WebhookHandler struct {
	notificationService ports.NotificationService
}

func NewWebhookHandler(notificationService ports.NotificationService) *WebhookHandler {
	return &WebhookHandler{
		notificationService: notificationService,
	}
}

func (h *WebhookHandler) HandleGitHubWebhook(c *gin.Context) {
	var payload map[string]interface{}
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(400, gin.H{"error": "Error al procesar el payload"})
		return
	}

	eventType := c.GetHeader("X-GitHub-Event")
	switch eventType {
	case "pull_request":
		pr := h.extractPullRequest(payload)
		if err := h.notificationService.HandlePullRequest(pr); err != nil {
			c.JSON(500, gin.H{"error": "Error al procesar el pull request"})
			return
		}
	case "workflow_run":
		workflow := h.extractWorkflow(payload)
		if err := h.notificationService.HandleWorkflow(workflow); err != nil {
			c.JSON(500, gin.H{"error": "Error al procesar el workflow"})
			return
		}
	case "workflow_dispatch":
		workflow := h.extractWorkflow(payload)
		if err := h.notificationService.HandleWorkflow(workflow); err != nil {
			c.JSON(500, gin.H{"error": "Error al procesar el workflow dispatch"})
			return
		}
	}

	c.JSON(200, gin.H{"status": "ok"})
}

func (h *WebhookHandler) HandleActionsWebhook(c *gin.Context) {
	var payload map[string]interface{}
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(400, gin.H{"error": "Error al procesar el payload"})
		return
	}

	action := h.extractAction(payload)
	if err := h.notificationService.HandleAction(action); err != nil {
		c.JSON(500, gin.H{"error": "Error al procesar la acción"})
		return
	}

	c.JSON(200, gin.H{"status": "ok"})
}

func (h *WebhookHandler) extractPullRequest(payload map[string]interface{}) domain.PullRequest {
	pr := payload["pull_request"].(map[string]interface{})
	user := pr["user"].(map[string]interface{})

	return domain.PullRequest{
		Action:   payload["action"].(string),
		Title:    pr["title"].(string),
		URL:      pr["html_url"].(string),
		Username: user["login"].(string),
		Merged:   pr["merged"].(bool),
	}
}

func (h *WebhookHandler) extractWorkflow(payload map[string]interface{}) domain.Workflow {
	workflow := payload["workflow"].(map[string]interface{})

	return domain.Workflow{
		Name:       workflow["name"].(string),
		Status:     payload["status"].(string),
		Conclusion: payload["conclusion"].(string),
	}
}

func (h *WebhookHandler) extractAction(payload map[string]interface{}) domain.Action {
	return domain.Action{
		Event:  payload["event"].(string),
		Status: payload["status"].(string),
	}
}
