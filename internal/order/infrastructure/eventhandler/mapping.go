package eventhandler

import (
	"github.com/google/uuid"

	"github.com/moeryomenko/saga/internal/order/domain"
	"github.com/moeryomenko/saga/schema"
)

func mapToDomainEvent(event map[string]any) (uuid.UUID, domain.Event, error) {
	switch kind := schema.GetEventType(event); kind {
	case schema.PaymentsConfirmed, schema.PaymentsFailed:
		return mapPaymentEventToDomain(schema.ToPaymentsEvent(event))
	case schema.StockConfirmed, schema.StockFailed:
		return mapStockEventToDomain(schema.ToStockEvent(event))
	default:
		panic(`bug: invalid event type`)
	}
}

func mapPaymentEventToDomain(event schema.PaymentsEvent, err error) (uuid.UUID, domain.Event, error) {
	if err != nil {
		return uuid.UUID{}, nil, err
	}

	switch event.Type {
	case schema.PaymentsConfirmed:
		return event.OrderID, domain.CofirmPayment{PaymentID: event.PaymentsID}, nil
	case schema.PaymentsFailed:
		return event.OrderID, domain.RejectPayment{}, nil
	}

	return uuid.UUID{}, nil, nil
}

func mapStockEventToDomain(event schema.StockEvent, err error) (uuid.UUID, domain.Event, error) {
	if err != nil {
		return uuid.UUID{}, nil, err
	}

	switch event.Type {
	case schema.StockConfirmed:
		return event.OrderID, domain.ConfirmStock{}, nil
	case schema.StockFailed:
		return event.OrderID, domain.RejectStock{}, nil
	}

	return uuid.UUID{}, nil, nil
}
