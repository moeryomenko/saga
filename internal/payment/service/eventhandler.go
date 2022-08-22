package service

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/cenkalti/backoff/v4"
	"github.com/google/uuid"
	"github.com/moeryomenko/saga/internal/payment/domain"
	"github.com/moeryomenko/saga/internal/payment/infrastructure/eventhandler"
	"github.com/moeryomenko/saga/internal/payment/infrastructure/repository"
)

func HandlePayments(ctx context.Context, customerID uuid.UUID, event domain.Event) error {
	_, err := repository.PersistTransaction(ctx, customerID, event)
	switch {
	case err == nil:
		return nil
	case errors.Is(err, domain.ErrDomain):
		return nil
	default:
		return err
	}
}

func Producer(pollPeriod time.Duration) func(ctx context.Context) error {
	return func(ctx context.Context) error {
		eventPollTicker := time.NewTicker(pollPeriod)
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
