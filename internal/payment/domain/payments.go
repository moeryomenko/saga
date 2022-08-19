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

type FailedPayment struct {
	ID      uuid.UUID
	OrderID uuid.UUID
	Amount  decimal.Decimal
}

func (p FailedPayment) GetID() uuid.UUID {
	return p.ID
}

func (p FailedPayment) GetAmount() decimal.Decimal {
	return p.Amount
}

func (p FailedPayment) GetOrderID() uuid.UUID {
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
	ID     uuid.UUID
	Amount decimal.Decimal
}

func (p CanceledPayment) GetID() uuid.UUID {
	return p.ID
}

func (p CanceledPayment) GetAmount() decimal.Decimal {
	return p.Amount
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

func FailPayment(payment Payment) FailedPayment {
	switch payment := payment.(type) {
	case NewPayment:
		return FailedPayment{
			ID:      payment.ID,
			OrderID: payment.OrderID,
			Amount:  payment.Amount,
		}
	default:
		panic(`bug: invalid failing payment flow`)
	}
}

func CancelPayment(payment Payment) (Payment, error) {
	switch payment := payment.(type) {
	case NewPayment, CompletedPayment:
		return CanceledPayment{
			ID:     payment.GetID(),
			Amount: payment.GetAmount(),
		}, nil
	default:
		return nil, ErrFailedPayment
	}
}
