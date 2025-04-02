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
		message = fmt.Sprintf("🆕 Nuevo Pull Request creado por %s\n**%s**\n%s", pr.Username, pr.Title, pr.URL)
	case "reopened":
		message = fmt.Sprintf("🔄 Pull Request reabierto por %s\n**%s**\n%s", pr.Username, pr.Title, pr.URL)
	case "ready_for_review":
		message = fmt.Sprintf("👀 Pull Request listo para revisión\n**%s**\n%s", pr.Title, pr.URL)
	case "closed":
		if pr.Merged {
			message = fmt.Sprintf("✅ Pull Request fusionado exitosamente\n**%s**\n%s", pr.Title, pr.URL)
		}
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
		switch workflow.Conclusion {
		case "success":
			emoji = "✅"
		case "failure":
			emoji = "❌"
		default:
			emoji = "⚠️"
		}

		message := fmt.Sprintf("%s Workflow **%s** %s\nEstado: %s",
			emoji,
			workflow.Name,
			workflow.Conclusion,
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
	message := fmt.Sprintf("🔄 GitHub Action ejecutado\nEvento: %s\nEstado: %s",
		action.Event,
		action.Status)

	notification := domain.Notification{
		ChannelID: s.channelRepo.GetTestsChannelID(),
		Message:   message,
		Type:      domain.ActionType,
	}
	return s.notificationRepo.Send(notification)
}
