package repository

import "CoolUrlShortener/internal/domain"

//go:generate go run github.com/vektra/mockery/v2@v2.42.1 --name EventsWriter
type EventsWriter interface {
	Insert(events []domain.URLEvent) error
}
