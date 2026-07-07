package sqlite

type VocabEventType string

const (
	VocabEventTypeAdd    VocabEventType = "add"
	VocabEventTypeRemove VocabEventType = "remove"
)

type VocabEvent struct {
	ID        string         `json:"id"`
	Type      VocabEventType `json:"type"`
	Word      string         `json:"word"`
	Timestamp int64          `json:"timestamp"`
}
