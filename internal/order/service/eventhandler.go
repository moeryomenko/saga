package service

import (
	"context"

	"github.com/google/uuid"
	"github.com/moeryomenko/saga/internal/order/domain"
	"github.com/moeryomenko/saga/internal/order/infrastructure/repository"
)

func HandleEvent(ctx context.Context, orderID uuid.UUID, event domain.Event) (domain.Order, error) {
	switch event.(type) {
	case domain.RejectPayment:
		// TODO: rollback saga transaction.
	case domain.RejectStock:
		// TODO rollback saga transaction.
	case domain.Process:
		// TODO: add beginning of saga transaction.
		// Send event to stream.
		return repository.PersistOrder(ctx, orderID, event)
	default:
		return repository.PersistOrder(ctx, orderID, event)
	}

	return nil, nil
}
