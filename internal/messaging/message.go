package messaging

const (
	// TimeCriticalSMS represents a necessary SMS message type
	TimeCriticalSMS = "timeCriticalSMS"
)

// Message represents any message that has to be sent
type Message struct {
	ID       string
	Content  string
	Receiver string
	Type     string
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
