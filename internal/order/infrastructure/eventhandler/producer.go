package eventhandler

import (
	"context"

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
				Items:      order.Items,
			}.Map(),
		}).Result()
		return err
	default:
		panic(`bug: failed state of order`)
	}
}
