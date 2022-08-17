package eventhandler

import (
	"context"
	"strings"

	redis "github.com/go-redis/redis/v8"

	"github.com/moeryomenko/saga/internal/order/domain"
	"github.com/moeryomenko/saga/schema"
)

func ProduceOrder(ctx context.Context, order domain.Order) error {
	switch order := order.(type) {
	case domain.PendingOrder:
		_, err := client.XAdd(ctx, &redis.XAddArgs{
			Stream: OrderStream,
			Values: schema.OrderEvent{
				Event:      schema.Event{Type: schema.NewOrder},
				OrderID:    order.ID,
				CustomerID: order.CustomerID,
				Price:      order.Price,
				Items:      mapItemsForEvent(order.Items),
			}.Map(),
		}).Result()
		return err
	default:
		panic(`bug: failed state of order`)
	}
}

func mapItemsForEvent(items []string) string {
	return strings.Join(items, `, `)
}

func ProduceRollback(ctx context.Context, order domain.Order) error {
	switch order := order.(type) {
	case domain.CanceledOrder:
		_, err := client.XAdd(ctx, &redis.XAddArgs{
			Stream: OrderStream,
			Values: schema.OrderEvent{
				Event:   schema.Event{Type: schema.CancelOrder},
				OrderID: order.GetID(),
			},
		}).Result()
		return err
	default:
		panic(`bug: invalid state for rollback`)
	}
}

func ProduceComplete(ctx context.Context, order domain.Order) error {
	switch order := order.(type) {
	case domain.CompletedOrder:
		_, err := client.XAdd(ctx, &redis.XAddArgs{
			Stream: OrderStream,
			Values: schema.OrderEvent{
				Event:   schema.Event{Type: schema.CompleteOrder},
				OrderID: order.GetID(),
			},
		}).Result()
		return err
	default:
		panic(`bug: invalid state for rollback`)
	}
}
