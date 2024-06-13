package events

import (
	"encoding/json"
	"log"
	"log/slog"

	"CoolUrlShortener/internal/repository"
	"CoolUrlShortener/internal/repository/models"
	"github.com/IBM/sarama"
)

type kafkaEventProducer struct {
	logger         *slog.Logger
	eventsProducer sarama.SyncProducer
}

func NewKafkaEventProducer(
	logger *slog.Logger,
	addrs []string,
	kafkaCfg *sarama.Config,
	doneCh <-chan struct{},
) (repository.EventsProducer, error) {
	producer, err := sarama.NewSyncProducer(addrs, kafkaCfg)
	if err != nil {
		return nil, err
	}
	go func() {
		<-doneCh
		err := producer.Close()
		if err != nil {
			logger.Error(err.Error())
		}
	}()

	return &kafkaEventProducer{
		logger:         logger,
		eventsProducer: producer,
	}, nil
}

func (k *kafkaEventProducer) ProduceEvent(event models.URLEvent) {
	bytes, err := json.Marshal(event)
	if err != nil {
		k.logger.Error(err.Error())
		return
	}

	log.Println(bytes)
	msg := &sarama.ProducerMessage{
		Topic: "events",
		Key:   sarama.StringEncoder(event.ShortURL),
		Value: sarama.ByteEncoder(bytes),
	}

	_, _, err = k.eventsProducer.SendMessage(msg)
	if err != nil {
		k.logger.Error(err.Error())
	}
}
