package services

import (
	"fmt"

	"github.com/actions/internal/core/domain"
	"github.com/actions/internal/core/ports"
)

type notificationService struct {
	notificationRepo ports.NotificationRepository
	channelRepo      ports.ChannelRepository
}

// NewNotificationService crea una nueva instancia del servicio de notificaciones 
//pruebas-02
func NewNotificationService(notificationRepo ports.NotificationRepository, channelRepo ports.ChannelRepository) ports.NotificationService {
	return &notificationService{
		notificationRepo: notificationRepo,
		channelRepo:      channelRepo,
	}
}

func (s *notificationService) HandlePullRequest(pr domain.PullRequest) error {
	var message string
	switch pr.Action {
	case "opened":
		message = fmt.Sprintf("🆕 **Nuevo Pull Request**\nCreado por: %s\nTítulo: %s\nURL: %s",
			pr.Username, pr.Title, pr.URL)
	case "reopened":
		message = fmt.Sprintf("🔄 **Pull Request Reabierto**\nPor: %s\nTítulo: %s\nURL: %s",
			pr.Username, pr.Title, pr.URL)
	case "ready_for_review":
		message = fmt.Sprintf("👀 **Pull Request Listo para Revisión**\nTítulo: %s\nURL: %s",
			pr.Title, pr.URL)
	case "closed":
		if pr.Merged {
			message = fmt.Sprintf("✅ **Pull Request Fusionado Exitosamente**\nTítulo: %s\nURL: %s",
				pr.Title, pr.URL)
		} else {
			message = fmt.Sprintf("❌ **Pull Request Cerrado sin Fusionar**\nTítulo: %s\nURL: %s",
				pr.Title, pr.URL)
		}
	case "synchronize":
		message = fmt.Sprintf("🔄 **Pull Request Actualizado**\nTítulo: %s\nURL: %s",
			pr.Title, pr.URL)
	}

	if message != "" {
		notification := domain.Notification{
			ChannelID: s.channelRepo.GetDevelopmentChannelID(),
			Message:   message,
			Type:      domain.PullRequestType,
		}
		return s.notificationRepo.Send(notification)
	}
	return nil
}

func (s *notificationService) HandleWorkflow(workflow domain.Workflow) error {
	if workflow.Status == "completed" {
		var emoji string
		var statusText string
		switch workflow.Conclusion {
		case "success":
			emoji = "✅"
			statusText = "Exitoso"
		case "failure":
			emoji = "❌"
			statusText = "Fallido"
		case "cancelled":
			emoji = "⚠️"
			statusText = "Cancelado"
		default:
			emoji = "⚠️"
			statusText = workflow.Conclusion
		}

		message := fmt.Sprintf("%s **Workflow Ejecutado**\nNombre: %s\nEstado: %s\nResultado: %s",
			emoji,
			workflow.Name,
			workflow.Status,
			statusText)

		notification := domain.Notification{
			ChannelID: s.channelRepo.GetTestsChannelID(),
			Message:   message,
			Type:      domain.WorkflowType,
		}
		return s.notificationRepo.Send(notification)
	}

	// Notificar también cuando el workflow comienza
	if workflow.Status == "queued" || workflow.Status == "in_progress" {
		message := fmt.Sprintf("🔄 **Workflow en Progreso**\nNombre: %s\nEstado: %s",
			workflow.Name,
			workflow.Status)

		notification := domain.Notification{
			ChannelID: s.channelRepo.GetTestsChannelID(),
			Message:   message,
			Type:      domain.WorkflowType,
		}
		return s.notificationRepo.Send(notification)
	}

	return nil
}

func (s *notificationService) HandleAction(action domain.Action) error {
	var emoji string
	switch action.Status {
	case "success":
		emoji = "✅"
	case "failure":
		emoji = "❌"
	default:
		emoji = "⚠️"
	}

	message := fmt.Sprintf("%s **GitHub Action**\nEvento: %s\nEstado: %s",
		emoji,
		action.Event,
		action.Status)

	notification := domain.Notification{
		ChannelID: s.channelRepo.GetTestsChannelID(),
		Message:   message,
		Type:      domain.ActionType,
	}
	return s.notificationRepo.Send(notification)
}
