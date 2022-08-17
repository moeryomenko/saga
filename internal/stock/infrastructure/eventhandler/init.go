package eventhandler

import (
	"context"

	redis "github.com/go-redis/redis/v8"

	"github.com/moeryomenko/saga/internal/stock/config"
)

var client *redis.Client = nil

func Init(cfg *config.Config) func(context.Context) error {
	return func(ctx context.Context) error {
		client = redis.NewClient(&redis.Options{Addr: cfg.Stream.Addr()})
		return client.Ping(ctx).Err()
	}
}

func Close(_ context.Context) error {
	return client.Close()
}

func HandleEvents(_ context.Context) error {
	return nil
}
