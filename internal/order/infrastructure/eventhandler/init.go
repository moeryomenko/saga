package eventhandler

import (
	"context"

	redis "github.com/go-redis/redis/v8"

	"github.com/moeryomenko/saga/internal/order/config"
)

var client *redis.Client = nil

const (
	OrderStream   = `orders_stream`
	ConfirmStream = `confirmation_stream`
	OrderGroup    = `orders_group`
)

func Init(cfg *config.Config) func(context.Context) error {
	return func(ctx context.Context) error {
		client = redis.NewClient(&redis.Options{Addr: cfg.Stream.Addr()})
		return client.Ping(ctx).Err()
	}
}

func initStreams(ctx context.Context) error {
	info, err := client.XInfoGroups(ctx, ConfirmStream).Result()
	if err != nil {
		return err
	}

	for _, groupInfo := range info {
		if groupInfo.Name == OrderGroup {
			return nil
		}
	}

	_, err = client.XGroupCreate(ctx, ConfirmStream, OrderGroup, `0`).Result()
	if err != nil {
		return err
	}

	return nil
}

func Close(_ context.Context) error {
	return client.Close()
}
