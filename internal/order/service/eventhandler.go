package service

import (
	"context"
	"log"
	"time"

	"github.com/cenkalti/backoff/v4"
	"github.com/google/uuid"
	"github.com/moeryomenko/saga/internal/order/domain"
	"github.com/moeryomenko/saga/internal/order/infrastructure/eventhandler"
	"github.com/moeryomenko/saga/internal/order/infrastructure/repository"
)

func HandleEvent(ctx context.Context, orderID uuid.UUID, event domain.Event) (domain.Order, error) {
	return repository.PersistOrder(ctx, orderID, event)
}

func Procuder(period time.Duration) func(ctx context.Context) error {
	return func(ctx context.Context) error {
		eventPollTicker := time.NewTicker(period)
		defer eventPollTicker.Stop()

		for {
			select {
			case <-ctx.Done():
				return nil
			case <-eventPollTicker.C:
				id, order, err := repository.GetEvent(ctx)
				switch err {
				case nil:
				case repository.ErrNoEvents:
					continue
				default:
					log.Println(err)
					return err
				}

				err = backoff.Retry(func() error {
					return eventhandler.Produce(ctx, order)
				}, backoff.NewExponentialBackOff())
				if err != nil {
					return err
				}

				err = repository.Ack(ctx, id)
				if err != nil {
					return err
				}
			}
		}
	}
}
