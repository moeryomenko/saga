package eventhandler

import (
	"context"
	"log"
	"time"

	redis "github.com/go-redis/redis/v8"
	"github.com/google/uuid"

	"github.com/moeryomenko/saga/internal/order/config"
	"github.com/moeryomenko/saga/internal/order/domain"
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
		err := client.Ping(ctx).Err()
		if err != nil {
			return err
		}
		return initStreams(ctx)
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

type EventHandler func(context.Context, uuid.UUID, domain.Event) (domain.Order, error)

func HandleEvents(handler EventHandler) func(context.Context) error {
	return func(ctx context.Context) error {
		for {
			select {
			case <-ctx.Done():
				return nil
			default:
				events, err := client.XReadGroup(ctx, &redis.XReadGroupArgs{
					Group:    OrderGroup,
					Consumer: uuid.NewString(),
					Streams:  []string{ConfirmStream, `>`},
					Block:    0,
					Count:    1,
					NoAck:    false,
				}).Result()
				if err != nil {
					<-time.After(time.Second)
					continue
				}

				for _, msg := range events[0].Messages {
					orderID, event, err := mapToDomainEvent(msg.Values)
					if err != nil {
						log.Println(err)
						_, _ = client.XAck(ctx, ConfirmStream, OrderGroup, msg.ID).Result()
						continue
					}

					_, err = handler(ctx, orderID, event)
					if err != nil {
						log.Println(err)
					}

					_, err = client.XAck(ctx, ConfirmStream, OrderGroup, msg.ID).Result()
					if err != nil {
						log.Println(err)
					}
				}
			}
		}
	}
}
