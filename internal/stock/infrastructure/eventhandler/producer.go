package eventhandler

import (
	"context"

	"github.com/go-redis/redis/v8"
	"github.com/moeryomenko/saga/internal/stock/domain"
	"github.com/moeryomenko/saga/schema"
)

func ProcudeConfimation(ctx context.Context, stock domain.Stock) error {
	event := schema.StockEvent{OrderID: stock.GetOrderID()}

	switch stock.(type) {
	case domain.ActiveStock:
		event.SetType(schema.StockConfirmed)
	case domain.RejectedStock:
		event.SetType(schema.StockFailed)
	default:
		panic(`bug: invalied state for stock`)
	}

	_, err := client.XAdd(ctx, &redis.XAddArgs{
		Stream: ConfirmStream,
		Values: event.Map(),
	}).Result()
	return err
}
