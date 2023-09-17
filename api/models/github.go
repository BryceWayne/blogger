package models

type WebhookPayload struct {
	Action     string `json:"action"`
	Issue      Issue  `json:"issue"`
	Repository Repo   `json:"repository"`
	Sender     User   `json:"sender"`
}

type Issue struct {
	URL    string `json:"url"`
	Number int    `json:"number"`
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
