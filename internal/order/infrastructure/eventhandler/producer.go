package eventhandler

import (
	"context"

	redis "github.com/go-redis/redis/v8"

	"github.com/moeryomenko/saga/schema"
)

func Produce(ctx context.Context, event schema.OrderEvent) error {
	_, err := client.XAdd(ctx, &redis.XAddArgs{
		Stream: OrderStream,
		Values: event.Map(),
	}).Result()
	return err
}
