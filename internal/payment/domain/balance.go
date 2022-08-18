package domain

import (
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type Balance struct {
	CustomerID uuid.UUID
	Amount     decimal.Decimal
	Reserved   decimal.Decimal
}

func (b Balance) Apply(payment Payment) (Balance, error) {
	switch payment := payment.(type) {
	case NewPayment:
		return b.ReserveAmount(payment)
	case CompletedPayment:
		return b.CompletePayment(payment), nil
	case CanceledPayment:
		return b.Refund(payment), nil
	default:
		panic(`bug: invalid payment`)
	}
}

func (b Balance) ReserveAmount(payment Payment) (Balance, error) {
	available, reserved := b.Amount.Sub(payment.GetAmount()), b.Reserved.Add(payment.GetAmount())

	if available.LessThan(decimal.Zero) {
		return b, ErrInsufficientFunds
	}

	return Balance{
		Amount:   available,
		Reserved: reserved,
	}, nil
}

func (b Balance) CompletePayment(payment Payment) Balance {
	return Balance{
		Amount:   b.Amount,
		Reserved: b.Reserved.Sub(payment.GetAmount()),
	}
}

func (b Balance) Refund(payment Payment) Balance {
	amount := payment.GetAmount()
	return Balance{
		Amount:   b.Amount.Add(amount),
		Reserved: b.Reserved.Sub(amount),
	}
}
