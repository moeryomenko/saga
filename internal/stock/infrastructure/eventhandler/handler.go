package eventhandler

import (
	"context"
	"log"
	"strings"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
	"github.com/moeryomenko/saga/internal/stock/domain"
	"github.com/moeryomenko/saga/schema"
)

type EventHandler func(domain.Event) (domain.Stock, error)

func HandleEvents(eventHandler EventHandler) func(ctx context.Context) error {
	return func(ctx context.Context) error {
		for initConsumerGroup(ctx) != nil {
		}

		consumerID := uuid.New().String()
		for {
			select {
			case <-ctx.Done():
				return nil
			default:
				events, err := client.XReadGroup(ctx, &redis.XReadGroupArgs{
					Group:    StockGroup,
					Consumer: consumerID,
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
					err = handleMessage(ctx, msg, eventHandler)
					if err != nil {
						log.Println(err)
					}
					_, _ = client.XAck(ctx, OrderStream, StockGroup, msg.ID).Result()
				}
			}
		}
	}
}

func handleMessage(ctx context.Context, msg redis.XMessage, eventHandler EventHandler) error {
	event, err := schema.ToOrderEvent(msg.Values)
	if err != nil {
		return err
	}

	// skip completed and canceled orders.
	if event.Type == schema.CompleteOrder || event.Type == schema.CancelOrder {
		return nil
	}

	stock, err := eventHandler(domain.StockOrder{
		OrderID: event.OrderID,
		Items:   mapItemsFromEvent(event.Items),
	})
	if err != nil {
		return err
	}

	err = ProcudeConfimation(ctx, stock)
	if err != nil {
		return err
	}

	_, err = client.XAck(ctx, OrderStream, StockGroup, msg.ID).Result()
	if err != nil {
		return err
	}
	return nil
}

func mapItemsFromEvent(items string) []string {
	return strings.Split(items, `, `)
}
