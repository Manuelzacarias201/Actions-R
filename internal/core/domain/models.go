package domain

// Notification representa una notificación que será enviada a Discord
type Notification struct {
	ChannelID string
	Message   string
	Type      NotificationType
}

// NotificationType representa el tipo de notificación
type NotificationType string

const (
	PullRequestType NotificationType = "pull_request"
	WorkflowType    NotificationType = "workflow"
	ActionType      NotificationType = "action"
)

// PullRequest representa un pull request de GitHub
type PullRequest struct {
	Action   string
	Title    string
	URL      string
	Username string
	Merged   bool
}

// Workflow representa un workflow de GitHub Actions
type Workflow struct {
	Name       string
	Status     string
	Conclusion string
}

// Action representa una acción de GitHub Actions
type Action struct {
	Event  string
	Status string
}
