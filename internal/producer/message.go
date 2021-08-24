package producer

type ActionType int

const (
	Create ActionType = iota
	Update
	Remove
)

type Message struct {
	Type ActionType
	Body EventMessage
}

type EventMessage struct {
	Action    string
	Id        uint64
	Timestamp int64
}

func CreateMessage(actionType ActionType, eventMessage EventMessage) Message {
	return Message{
		Type: actionType,
		Body: eventMessage,
	}
}

func (actionType ActionType) String() string {
	switch actionType {
	case Create:
		return "Created"
	case Update:
		return "Updated"
	case Remove:
		return "Removed"
	default:
		return "Unknown MessageType"
	}
}