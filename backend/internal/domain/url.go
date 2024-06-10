package domain

import "time"

type URLData struct {
	ID        int64
	ShortUrl  string
	LongUrl   string
	CreatedAt time.Time
}

type URLEvent struct {
	LongURL   string
	ShortURL  string
	EventTime time.Time
	EventType int8
}
