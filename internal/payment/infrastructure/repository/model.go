package repository

import (
	"github.com/google/uuid"
	"github.com/jackc/pgtype"
	"github.com/moeryomenko/saga/internal/payment/domain"
	"github.com/shopspring/decimal"
)

const (
	statusNew       = `new`
	statusFailed    = `failed`
	statusCompleted = `completed`
	statusCanceled  = `canceled`
)

type Balance struct {
	CustomerID pgtype.UUID
	Available  decimal.Decimal
	Reserved   decimal.Decimal
}

type Payment struct {
	PaymentID  pgtype.UUID
	CustomerID pgtype.UUID
	OrderID    pgtype.UUID
	Amount     decimal.Decimal
	Status     string
}

func mapPaymentToDomain(p *Payment) domain.Payment {
	switch p.Status {
	case statusNew:
		return domain.NewPayment{
			ID:      p.PaymentID.Bytes,
			OrderID: p.OrderID.Bytes,
			Amount:  p.Amount,
		}
	case statusFailed:
		return domain.FailedPayment{
			ID:      p.PaymentID.Bytes,
			OrderID: p.OrderID.Bytes,
			Amount:  p.Amount,
		}
	case statusCompleted:
		return domain.CompletedPayment{
			ID:     p.PaymentID.Bytes,
			Amount: p.Amount,
		}
	case statusCanceled:
		return domain.CanceledPayment{
			ID:     p.PaymentID.Bytes,
			Amount: p.Amount,
		}
	}
	return nil
}

func mapBalanceToDomain(b *Balance) domain.Balance {
	return domain.Balance{
		CustomerID: b.CustomerID.Bytes,
		Amount:     b.Available,
		Reserved:   b.Reserved,
	}
}

func mapPaymentToModel(customerID uuid.UUID, p domain.Payment) Payment {
	status := statusNew
	switch p := p.(type) {
	case domain.NewPayment:
		return Payment{
			PaymentID:  pgtype.UUID{Bytes: p.ID, Status: pgtype.Present},
			OrderID:    pgtype.UUID{Bytes: p.OrderID, Status: pgtype.Present},
			CustomerID: pgtype.UUID{Bytes: customerID, Status: pgtype.Present},
			Amount:     p.Amount,
			Status:     status,
		}
	case domain.FailedPayment:
		return Payment{
			PaymentID:  pgtype.UUID{Bytes: p.ID, Status: pgtype.Present},
			OrderID:    pgtype.UUID{Bytes: p.OrderID, Status: pgtype.Present},
			CustomerID: pgtype.UUID{Bytes: customerID, Status: pgtype.Present},
			Amount:     p.Amount,
			Status:     statusFailed,
		}
	case domain.CompletedPayment:
		status = statusCompleted
	case domain.CanceledPayment:
		status = statusCanceled
	}
	return Payment{
		PaymentID: pgtype.UUID{Bytes: p.GetID(), Status: pgtype.Present},
		Status:    status,
	}
}

func mapBalanceToModel(b domain.Balance) Balance {
	return Balance{
		CustomerID: pgtype.UUID{Bytes: b.CustomerID, Status: pgtype.Present},
		Available:  b.Amount,
		Reserved:   b.Reserved,
	}
}
