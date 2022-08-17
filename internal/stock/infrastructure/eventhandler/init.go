package eventhandler

import (
	"context"

	redis "github.com/go-redis/redis/v8"

	"github.com/moeryomenko/saga/internal/stock/config"
)

var client *redis.Client = nil

const (
	OrderStream    = `orders_stream`
	RollbackStream = `rollback_strem`
	StockGroup     = `stock_group`
)

func Init(cfg *config.Config) func(context.Context) error {
	return func(ctx context.Context) error {
		client = redis.NewClient(&redis.Options{Addr: cfg.Stream.Addr()})
		err := client.Ping(ctx).Err()
		if err != nil {
			return err
		}
		return initConsumerGroup(ctx)
	}
}

func initConsumerGroup(ctx context.Context) error {
	info, err := client.XInfoGroups(ctx, OrderStream).Result()
	if err != nil {
		return err
	}

	for _, groupInfo := range info {
		if groupInfo.Name == StockGroup {
			return nil
		}
	}

	_, err = client.XGroupCreate(ctx, OrderStream, StockGroup, `$`).Result()
	if err != nil {
		return err
	}

	return nil
}

func Close(_ context.Context) error {
	return client.Close()
}

func HandleEvents(_ context.Context) error {
	return nil
}
