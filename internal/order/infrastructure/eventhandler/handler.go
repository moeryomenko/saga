package eventhandler

import (
	"context"
	"log"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"

	"github.com/moeryomenko/saga/internal/order/domain"
)

type EventHandler func(context.Context, uuid.UUID, domain.Event) (domain.Order, error)

func HandleEvents(handler EventHandler) func(context.Context) error {
	return func(ctx context.Context) error {
		for initStreams(ctx) != nil {
		}

		for {
			select {
			case <-ctx.Done():
				return nil
			default:
				events, err := client.XReadGroup(ctx, &redis.XReadGroupArgs{
					Group:    OrderGroup,
					Consumer: uuid.NewString(),
					Streams:  []string{ConfirmStream, `>`},
					Block:    0,
					Count:    1,
					NoAck:    false,
				}).Result()
				if err != nil {
					<-time.After(time.Second)
					continue
				}

				for _, msg := range events[0].Messages {
					err = handleMessage(ctx, msg, handler)
					if err != nil {
						log.Println(err)
					}

					_, err = client.XAck(ctx, ConfirmStream, OrderGroup, msg.ID).Result()
					if err != nil {
						log.Println(err)
					}
				}
			}
		}
	}
}

func handleMessage(ctx context.Context, msg redis.XMessage, eventHandler EventHandler) error {
	orderID, event, err := mapToDomainEvent(msg.Values)
	if err != nil {
		return err
	}

	_, err = eventHandler(ctx, orderID, event)
	return err
}
