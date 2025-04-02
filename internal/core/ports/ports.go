package ports

import "github.com/actions/internal/core/domain"

// NotificationService define el puerto primario para el servicio de notificaciones
type NotificationService interface {
	HandlePullRequest(pr domain.PullRequest) error
	HandleWorkflow(workflow domain.Workflow) error
	HandleAction(action domain.Action) error
}

// NotificationRepository define el puerto secundario para enviar notificaciones
type NotificationRepository interface {
	Send(notification domain.Notification) error
	Initialize() error
}

// ChannelRepository define el puerto secundario para gestionar canales
type ChannelRepository interface {
	CreateChannels() error
	GetDevelopmentChannelID() string
	GetTestsChannelID() string
}
