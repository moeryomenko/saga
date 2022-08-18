package domain

import (
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type Payment interface {
	GetID() uuid.UUID
	GetAmount() decimal.Decimal
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

func CancelPayment(payment Payment) (Payment, error) {
	switch payment := payment.(type) {
	case NewPayment:
		return CanceledPayment{
			ID:     payment.ID,
			Amount: payment.Amount,
		}, nil
	default:
		return nil, ErrCompletedPayment
	}
}
