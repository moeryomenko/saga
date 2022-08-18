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

func (b Balance) Apply(payment Payment) (balance Balance, err error) {
	switch payment := payment.(type) {
	case NewPayment:
		balance, err = b.ReserveAmount(payment)
		if err != nil {
			return b, err
		}
	case CompletedPayment:
		balance = b.CompletePayment(payment)
	case CanceledPayment:
		balance = b.Refund(payment)
	default:
		panic(`bug: invalid payment`)
	}
	balance.CustomerID = b.CustomerID
	return balance, nil
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
