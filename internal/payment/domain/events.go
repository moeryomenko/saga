package domain

import (
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type Event interface {
	GetID() uuid.UUID
}

type Reserve struct {
	OrderID uuid.UUID
	Amount  decimal.Decimal
}

func (e Reserve) GetID() uuid.UUID {
	return uuid.UUID{}
}

type Complete struct {
	PaymentID uuid.UUID
}

func (e Complete) GetID() uuid.UUID {
	return e.PaymentID
}

type Cancel struct {
	PaymentID uuid.UUID
}

func (e Cancel) GetID() uuid.UUID {
	return e.PaymentID
}

func Apply(payment Payment, event Event) (Payment, error) {
	switch event := event.(type) {
	case Reserve:
		return NewPayment{
			ID:      uuid.New(),
			OrderID: event.OrderID,
			Amount:  event.Amount,
		}, nil
	case Complete:
		return CompletePayment(payment)
	case Cancel:
		return CancelPayment(payment)
	default:
		panic(`bug: invalid payment event`)
	}
}
