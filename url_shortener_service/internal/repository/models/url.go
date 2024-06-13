package models

const (
	EventTypeCreate = 1
	EventTypeFollow = 2
)

type URLEvent struct {
	LongURL   string `json:"long_url"`
	ShortURL  string `json:"short_url"`
	EventTime int64  `json:"event_time"`
	EventType int8   `json:"event_type"`
}
