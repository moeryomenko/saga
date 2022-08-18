package service

import (
	"context"

	"github.com/google/uuid"
	"github.com/moeryomenko/saga/internal/payment/domain"
	"github.com/moeryomenko/saga/internal/payment/infrastructure/eventhandler"
	"github.com/moeryomenko/saga/internal/payment/infrastructure/repository"
)

func HandlePayments(ctx context.Context, customerID uuid.UUID, event domain.Event) error {
	payment, err := repository.PersistTransaction(ctx, customerID, event)
	if err != nil {
		return err
	}
	switch payment := payment.(type) {
	case domain.ResultPayment:
		return eventhandler.Produce(ctx, payment.GetOrderID(), payment)
	case domain.CompletedPayment:
		return nil
	default:
		panic(`bug: invalid payment result`)
	}
}
