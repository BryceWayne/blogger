package models

type WebhookPayload struct {
	Action     string    `json:"action"`
	Issue      Issue     `json:"issue"`
	Repository Repo      `json:"repository"`
	Sender     User      `json:"sender"`
	Commits    []Commits `json:"commits"`
}

type Issue struct {
	URL    string `json:"url"`
	Number int    `json:"number"`
	// Add other fields as needed
}

type Author struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Username string `json:"username"`
	// Add other fields as needed
}

type Repo struct {
	ID       int    `json:"id"`
	FullName string `json:"full_name"`
	Owner    User   `json:"owner"`
	// Add other fields as needed
}

type User struct {
	Login string `json:"login"`
	ID    int    `json:"id"`
	// Add other fields as needed
}

type Commits struct {
	Added    []string `json:"added"`
	Removed  []string `json:"removed"`
	Modified []string `json:"modified"`
}

type FileStatus struct {
	Added    bool
	Modified bool
	Removed  bool
}
