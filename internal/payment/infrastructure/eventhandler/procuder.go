package eventhandler

import (
	"context"

	"github.com/go-redis/redis/v8"
	"github.com/moeryomenko/saga/schema"
)

func Produce(ctx context.Context, event schema.PaymentsEvent) error {
	_, err := client.XAdd(ctx, &redis.XAddArgs{
		Stream: ConfirmStream,
		Values: event.Map(),
	}).Result()
	return err
}
