package domain

import (
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type Payment interface {
	GetID() uuid.UUID
	GetAmount() decimal.Decimal
}

type ResultPayment interface {
	Payment
	GetOrderID() uuid.UUID
}

type NewPayment struct {
	ID      uuid.UUID
	OrderID uuid.UUID
	Amount  decimal.Decimal
}

func (p NewPayment) GetID() uuid.UUID {
	return p.ID
}

func (p NewPayment) GetAmount() decimal.Decimal {
	return p.Amount
}

func (p NewPayment) GetOrderID() uuid.UUID {
	return p.OrderID
}

type CompletedPayment struct {
	ID     uuid.UUID
	Amount decimal.Decimal
}

func (p CompletedPayment) GetID() uuid.UUID {
	return p.ID
}

func (p CompletedPayment) GetAmount() decimal.Decimal {
	return p.Amount
}

type CanceledPayment struct {
	OrderID uuid.UUID
	ID      uuid.UUID
	Amount  decimal.Decimal
}

func (p CanceledPayment) GetID() uuid.UUID {
	return p.ID
}

func (p CanceledPayment) GetAmount() decimal.Decimal {
	return p.Amount
}

func (p CanceledPayment) GetOrderID() uuid.UUID {
	return p.OrderID
}

func CompletePayment(payment Payment) (Payment, error) {
	switch payment := payment.(type) {
	case NewPayment:
		return CompletedPayment{
			ID:     payment.ID,
			Amount: payment.Amount,
		}, nil
	default:
		return nil, ErrCanceledPayment
	}
}

func CancelPayment(payment Payment) (Payment, error) {
	switch payment := payment.(type) {
	case NewPayment:
		return CanceledPayment{
			OrderID: payment.OrderID,
			ID:      payment.ID,
			Amount:  payment.Amount,
		}, nil
	default:
		return nil, ErrCompletedPayment
	}
}
