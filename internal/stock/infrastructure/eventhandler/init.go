package eventhandler

import (
	"context"
	"log"
	"strings"
	"time"

	redis "github.com/go-redis/redis/v8"
	"github.com/google/uuid"

	"github.com/moeryomenko/saga/internal/stock/config"
	"github.com/moeryomenko/saga/internal/stock/domain"
	"github.com/moeryomenko/saga/schema"
)

var client *redis.Client = nil

const (
	OrderStream   = `orders_stream`
	ConfirmStream = `confirmation_stream`
	StockGroup    = `stock_group`
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
		if groupInfo.Name == StockGroup {
			return nil
		}
	}

	_, err = client.XGroupCreate(ctx, OrderStream, StockGroup, `0`).Result()
	if err != nil {
		return err
	}

	return nil
}

func Close(_ context.Context) error {
	return client.Close()
}

type EventHandler func(domain.Event) (domain.Stock, error)

func HandleEvents(eventHandler EventHandler) func(ctx context.Context) error {
	return func(ctx context.Context) error {
		for initConsumerGroup(ctx) != nil {
		}

		for {
			select {
			case <-ctx.Done():
				return nil
			default:
				events, err := client.XReadGroup(ctx, &redis.XReadGroupArgs{
					Group:    StockGroup,
					Consumer: uuid.NewString(),
					Streams:  []string{OrderStream, `>`},
					Block:    0,
					Count:    1,
					NoAck:    false,
				}).Result()
				if err != nil {
					<-time.After(time.Second)
					continue
				}

				for _, msg := range events[0].Messages {
					event, err := schema.ToOrderEvent(msg.Values)
					if err != nil {
						log.Println(err)
						_, _ = client.XAck(ctx, OrderStream, StockGroup, msg.ID).Result()
						continue
					}

					stock, err := eventHandler(domain.StockOrder{
						OrderID: event.OrderID,
						Items:   mapItemsFromEvent(event.Items),
					})
					if err != nil {
						log.Println(err)
						_, _ = client.XAck(ctx, OrderStream, StockGroup, msg.ID).Result()
						continue
					}

					err = ProcudeConfimation(ctx, stock)
					if err != nil {
						log.Println(err)
					}

					_, err = client.XAck(ctx, OrderStream, StockGroup, msg.ID).Result()
					if err != nil {
						log.Println(err)
					}
				}
			}
		}
	}
}

func mapItemsFromEvent(items string) []string {
	return strings.Split(items, `, `)
}

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
