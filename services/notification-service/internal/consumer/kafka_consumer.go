package consumer

import (
	"context"
	"encoding/json"
	"log/slog"

	"github.com/maxluvds/booking-platform/pkg/kafka"
	"github.com/maxluvds/booking-platform/services/notification-service/internal/usecase"
)

type KafkaConsumer struct {
	consumer *kafka.Consumer
	usecase  *usecase.NotificationUseCase
	log      *slog.Logger
}

func NewKafkaConsumer(
	consumer *kafka.Consumer,
	usecase *usecase.NotificationUseCase,
	log *slog.Logger,
) *KafkaConsumer {
	return &KafkaConsumer{
		consumer: consumer,
		usecase:  usecase,
		log:      log,
	}
}

func (c *KafkaConsumer) Start(ctx context.Context) {
	c.log.Info("Kafka consumer started")

	for {
		select {
		case <-ctx.Done():
			c.log.Info("Kafka consumer stopping...")
			return
		default:
			msg, err := c.consumer.Read(ctx)
			if err != nil {
				c.log.Error("Failed to read message", "error", err)
				continue
			}

			c.log.Info("Received message",
				"topic", msg.Topic,
				"key", string(msg.Key),
			)

			var eventMsg kafka.EventCreatedMessage
			if err := json.Unmarshal(msg.Value, &eventMsg); err != nil {
				c.log.Error("Failed to unmarshal message", "error", err)
				continue
			}

			if err := c.usecase.HandleEventCreated(ctx, eventMsg); err != nil {
				c.log.Error("Failed to handle event", "error", err)
				continue
			}

			c.log.Info("Notification saved",
				"event_id", eventMsg.EventID,
				"title", eventMsg.Title,
			)
		}
	}
}
