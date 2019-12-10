package messaging

// Message represents any message that has to be sent
type Message struct {
	ID        string `json:"id"`
	Critical  bool   `json:"critical"`
	Recipient string `json:"recipient"`
	Content   string `json:"content"`
}

// MessageResource represents a sent message
type MessageResource struct {
	Message
	Status     string `json:"status"`
	ProviderID string `json:"provider_id"`
}

// Sender is a provider which is able to send messages
type Sender interface {
	Send(message Message) (MessageResource, error)
}
