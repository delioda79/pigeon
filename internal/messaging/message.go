package messaging

// Message represents any message that has to be sent
type Message struct {
	ID        string
	Critical  bool
	Recipient string
	Content   string
}

// MessageResource represents a sent message
type MessageResource struct {
	Message
	Status     string
	ProviderID string
}

// Sender is a provider which is able to send messages
type Sender interface {
	Send(message Message) (MessageResource, error)
}
