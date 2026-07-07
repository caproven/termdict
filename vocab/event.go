package vocab

type EventType string

const (
	EventTypeAdd    EventType = "add"
	EventTypeRemove EventType = "remove"
)

type Event struct {
	ID        string    `json:"id"`
	Type      EventType `json:"type"`
	Word      string    `json:"word"`
	Timestamp int64     `json:"timestamp"`
}
