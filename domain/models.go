package domain

type PullRequestEvent struct {
	Action      string      `json:"action"`
	Number      int         `json:"number"`
	PullRequest PullRequest `json:"pull_request"`
	Repository  Repository  `json:"repository"`
}

type PullRequest struct {
	Title   string `json:"title"`
	HTMLURL string `json:"html_url"`
	State   string `json:"state"`
	User    User   `json:"user"`
	Body    string `json:"body"`
	Base    Base   `json:"base"`
	Head    Head   `json:"head"`
	Merged  bool   `json:"merged"`
}

type Repository struct {
	FullName string `json:"full_name"`
	HTMLURL  string `json:"html_url"`
}

type User struct {
	Login     string `json:"login"`
	AvatarURL string `json:"avatar_url"`
}

type Base struct {
	Ref string `json:"ref"`
}

type Head struct {
	Ref string `json:"ref"`
}

type WorkflowRunEvent struct {
	Action     string     `json:"action"`
	Workflow   Workflow   `json:"workflow"`
	Repository Repository `json:"repository"`
	Status     string     `json:"status"`
	Conclusion string     `json:"conclusion"`
	HTMLURL    string     `json:"html_url"`
	CreatedAt  string     `json:"created_at"`
	UpdatedAt  string     `json:"updated_at"`
	HeadBranch string     `json:"head_branch"`
	//HeadCommit      Commit     `json:"head_commit"`
	Event           string `json:"event"`
	RunAttempt      int    `json:"run_attempt"`
	RunStartedAt    string `json:"run_started_at"`
	TriggeringActor User   `json:"triggering_actor"`
}

type Workflow struct {
	Name string `json:"name"`
	ID   int    `json:"id"`
}

type WorkflowJobEvent struct {
	Action      string      `json:"action"`
	WorkflowJob WorkflowJob `json:"workflow_job"`
	Repository  Repository  `json:"repository"`
	Sender      User        `json:"sender"`
}

type WorkflowJob struct {
	ID          int64          `json:"id"`
	RunID       int64          `json:"run_id"`
	RunURL      string         `json:"run_url"`
	NodeID      string         `json:"node_id"`
	HeadSha     string         `json:"head_sha"`
	URL         string         `json:"url"`
	HTMLURL     string         `json:"html_url"`
	Status      string         `json:"status"`
	Conclusion  string         `json:"conclusion"`
	StartedAt   string         `json:"started_at"`
	CompletedAt string         `json:"completed_at"`
	Name        string         `json:"name"`
	Steps       []WorkflowStep `json:"steps"`
}

type WorkflowStep struct {
	Name        string `json:"name"`
	Status      string `json:"status"`
	Conclusion  string `json:"conclusion"`
	Number      int    `json:"number"`
	StartedAt   string `json:"started_at"`
	CompletedAt string `json:"completed_at"`
}
