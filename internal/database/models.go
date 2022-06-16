package database

type Message struct {
	ID      string   `json:"id" binding:"required"`
	Body    string   `json:"body" binding:"required"`
	Actions []Action `json:"actions,omitempty"`
	Status  string   `json:"status"`
	Reason  string   `json:"reasons,omitempty"`
}

type Action struct {
	ID     string `json:"id"`
	Status string `json:"status"`
	Reason string `json:"reason"`
}

type Approval struct {
	ID        string `json:"id"`
	Status    string `json:"status"`
	MessageID string `json:"messageId"`
	Reason    string `json:"reason"`
}
