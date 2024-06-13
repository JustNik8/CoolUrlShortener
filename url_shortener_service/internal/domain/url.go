package domain

import "time"

type URLData struct {
	ID        int64
	ShortUrl  string
	LongUrl   string
	CreatedAt time.Time
}
