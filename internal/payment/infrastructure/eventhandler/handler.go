package eventhandler

import (
	"context"
	"log"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"

	"github.com/moeryomenko/saga/internal/payment/domain"
	"github.com/moeryomenko/saga/schema"
)

type EventHandler func(context.Context, uuid.UUID, domain.Event) error

func HandleEvents(handler EventHandler) func(ctx context.Context) error {
	return func(ctx context.Context) error {
		for initConsumerGroup(ctx) != nil {
		}

		for {
			select {
			case <-ctx.Done():
				return nil
			default:
				events, err := client.XReadGroup(ctx, &redis.XReadGroupArgs{
					Group:    PaymentGroup,
					Consumer: uuid.NewString(),
					Streams:  []string{OrderStream, `>`},
					Block:    0,
					Count:    1,
					NoAck:    false,
				}).Result()
				if err != nil {
					<-time.After(time.Second)
					continue
				}

				for _, msg := range events[0].Messages {
					err := handleMessage(ctx, msg, handler)
					if err != nil {
						log.Panicln(err)
					}

					_, err = client.XAck(ctx, OrderStream, PaymentGroup, msg.ID).Result()
					if err != nil {
						log.Println(err)
					}
				}
			}
		}
	}
}

func handleMessage(ctx context.Context, msg redis.XMessage, handler EventHandler) error {
	event, err := schema.ToOrderEvent(msg.Values)
	if err != nil {
		return err
	}

	var domainEvent domain.Event
	switch event.Type {
	case schema.NewOrder:
		domainEvent = domain.Reserve{OrderID: event.OrderID, Amount: event.Price}
	case schema.CompleteOrder:
		domainEvent = domain.Complete{PaymentID: event.PaymentID}
	case schema.CancelOrder:
		domainEvent = domain.Cancel{PaymentID: event.PaymentID}
	}

	return handler(ctx, event.CustomerID, domainEvent)
}
