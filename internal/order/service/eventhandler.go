package service

import (
	"context"

	"github.com/google/uuid"
	"github.com/moeryomenko/saga/internal/order/domain"
	"github.com/moeryomenko/saga/internal/order/infrastructure/eventhandler"
	"github.com/moeryomenko/saga/internal/order/infrastructure/repository"
)

func HandleEvent(ctx context.Context, orderID uuid.UUID, event domain.Event) (domain.Order, error) {
	switch event.(type) {
	case domain.CofirmPayment, domain.ConfirmStock:
		order, err := repository.PersistOrder(ctx, orderID, event)
		if err != nil {
			return nil, err
		}

		if order, ok := order.(domain.CompletedOrder); ok {
			err = eventhandler.Produce(ctx, order)
		}
		return order, err
	case domain.RejectPayment, domain.RejectStock:
		order, err := repository.PersistOrder(ctx, orderID, event)
		if err != nil {
			return nil, err
		}

		return order, eventhandler.Produce(ctx, order)
	case domain.Process:
		order, err := repository.PersistOrder(ctx, orderID, event)
		if err != nil {
			return nil, err
		}

		return order, eventhandler.Produce(ctx, order)
	default:
		return repository.PersistOrder(ctx, orderID, event)
	}
}
