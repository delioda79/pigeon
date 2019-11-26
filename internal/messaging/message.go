package messaging

const (
	Necessary = "necessary"
)

type Message struct {
	ID       string
	Content  string
	Receiver string
	Type     string
}

type Sender interface {
	Send(message Message) error
}
