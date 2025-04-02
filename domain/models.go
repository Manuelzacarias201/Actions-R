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
}

type Workflow struct {
	Name string `json:"name"`
	ID   int    `json:"id"`
}
