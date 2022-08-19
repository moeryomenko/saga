package service

import (
	"context"

	"github.com/google/uuid"
	"github.com/moeryomenko/saga/internal/order/domain"
	"github.com/moeryomenko/saga/internal/order/infrastructure/eventhandler"
	"github.com/moeryomenko/saga/internal/order/infrastructure/repository"
)

func HandleEvent(ctx context.Context, orderID uuid.UUID, event domain.Event) (domain.Order, error) {
	order, err := repository.PersistOrder(ctx, orderID, event)
	if err != nil {
		return nil, err
	}
	switch event.(type) {
	case domain.CofirmPayment, domain.ConfirmStock,
		domain.RejectPayment, domain.RejectStock,
		domain.Process:
		return order, eventhandler.Produce(ctx, order)
	}

	return order, nil
}
