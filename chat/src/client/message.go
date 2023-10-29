package client

type Message struct {
	Sender    string `json:"sender,omitempty"`
	Recipient string `json:"recipient,omitempty"`
	Content   string `json:"content,omitempty"`
	// Unique ID for the Message created.
	Id        string `json:"id"`
	Timestamp int64  `json:"timestamp"`
}
