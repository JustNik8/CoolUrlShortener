package repository

import (
	"CoolUrlShortener/internal/repository/models"
)

//go:generate go run github.com/vektra/mockery/v2@v2.42.1 --name EventsServiceProducer
type EventsProducer interface {
	ProduceEvent(event models.URLEvent)
}
