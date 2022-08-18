package eventhandler

import (
	"context"
	"strings"

	redis "github.com/go-redis/redis/v8"

	"github.com/moeryomenko/saga/internal/order/domain"
	"github.com/moeryomenko/saga/schema"
)

func Produce(ctx context.Context, order domain.Order) error {
	var event any
	switch order := order.(type) {
	case domain.PendingOrder:
		event = schema.OrderEvent{
			Event:      schema.Event{Type: schema.NewOrder},
			OrderID:    order.ID,
			CustomerID: order.CustomerID,
			Price:      order.Price,
			Items:      mapItemsForEvent(order.Items),
		}.Map()
	case domain.CanceledOrder:
		event = schema.OrderEvent{
			Event:   schema.Event{Type: schema.CancelOrder},
			OrderID: order.GetID(),
		}.Map()
	case domain.CompletedOrder:
		event = schema.OrderEvent{
			Event:   schema.Event{Type: schema.CompleteOrder},
			OrderID: order.GetID(),
		}.Map()
	default:
		panic(`bug: failed state of order`)
	}

	_, err := client.XAdd(ctx, &redis.XAddArgs{
		Stream: OrderStream,
		Values: event,
	}).Result()
	return err
}

func mapItemsForEvent(items []string) string {
	return strings.Join(items, `, `)
}
