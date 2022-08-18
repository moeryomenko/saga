package eventhandler

import (
	"context"

	redis "github.com/go-redis/redis/v8"

	"github.com/moeryomenko/saga/internal/payment/config"
)

var client *redis.Client = nil

const (
	OrderStream   = `orders_stream`
	ConfirmStream = `confirmation_stream`
	PaymentGroup  = `payments_group`
)

func Init(cfg *config.Config) func(context.Context) error {
	return func(ctx context.Context) error {
		client = redis.NewClient(&redis.Options{Addr: cfg.Stream.Addr()})
		return client.Ping(ctx).Err()
	}
}

func initConsumerGroup(ctx context.Context) error {
	info, err := client.XInfoGroups(ctx, OrderStream).Result()
	if err != nil {
		return err
	}

	for _, groupInfo := range info {
		if groupInfo.Name == PaymentGroup {
			return nil
		}
	}

	_, err = client.XGroupCreate(ctx, OrderStream, PaymentGroup, `0`).Result()
	if err != nil {
		return err
	}

	return nil
}

func Close(_ context.Context) error {
	return client.Close()
}
