package repository

import "CoolUrlShortener/internal/domain"

type EventsWriter interface {
	Insert(events []domain.URLEvent) error
}
