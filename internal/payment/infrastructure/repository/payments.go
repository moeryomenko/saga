package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgtype"
	"github.com/jackc/pgx/v4"
	"github.com/moeryomenko/saga/internal/payment/domain"
	"github.com/moeryomenko/saga/pkg/errors"
)

func PersistTransaction(ctx context.Context, customerID uuid.UUID, event domain.Event) (domain.Payment, error) {
	var payment domain.Payment
	err := pool.BeginTxFunc(ctx, pgx.TxOptions{IsoLevel: pgx.ReadCommitted}, func(tx pgx.Tx) error {
		balance, err := findBalanceByCustomer(ctx, tx, customerID)
		if err != nil {
			return errors.MarkAndWrapError(err, ErrInfrastructure, `couldn't find balance`)
		}

		payment, err = findPaymentByID(ctx, tx, event.GetID())
		switch err {
		case nil, pgx.ErrNoRows:
			// it's ok for NewPayment event.
		default:
			return errors.MarkAndWrapError(err, ErrInfrastructure, `couldn't find payment`)
		}

		balance, payment, err = balance.Transaction(domain.Tx{
			Payment: payment,
			Event:   event,
		})
		if err != nil {
			return errors.MarkAndWrapError(err, domain.ErrDomain, `couldn't apply event`)
		}

		err = saveBalance(ctx, tx, balance)
		if err != nil {
			return err
		}

		return savePayment(ctx, tx, customerID, payment)
	})
	return payment, err
}

func findBalanceByCustomer(ctx context.Context, tx pgx.Tx, customerID uuid.UUID) (domain.Balance, error) {
	balance := &Balance{CustomerID: pgtype.UUID{Bytes: customerID, Status: pgtype.Present}}
	err := tx.QueryRow(ctx, findBalanceQuery, customerID.String()).Scan(
		&balance.Available,
		&balance.Reserved,
	)
	if err != nil {
		return domain.Balance{}, err
	}
	return mapBalanceToDomain(balance), nil
}

func saveBalance(ctx context.Context, tx pgx.Tx, balance domain.Balance) error {
	model := mapBalanceToModel(balance)
	_, err := tx.Exec(ctx, updateBalanceQuery, model.CustomerID, model.Available, model.Reserved)
	return err
}

func findPaymentByID(ctx context.Context, tx pgx.Tx, paymentID uuid.UUID) (domain.Payment, error) {
	payment := &Payment{PaymentID: pgtype.UUID{Bytes: paymentID, Status: pgtype.Present}}
	err := tx.QueryRow(ctx, findPaymentQuery, paymentID.String()).Scan(
		&payment.CustomerID,
		&payment.OrderID,
		&payment.Amount,
		&payment.Status,
	)
	if err != nil {
		return nil, err
	}
	return mapPaymentToDomain(payment), nil
}

func savePayment(ctx context.Context, tx pgx.Tx, customerID uuid.UUID, payment domain.Payment) error {
	model := mapPaymentToModel(customerID, payment)
	switch model.Status {
	case statusNew:
		_, err := tx.Exec(ctx, insertPaymentQuery, model.PaymentID, model.Status, model.CustomerID, model.OrderID, model.Amount)
		return err
	case statusFailed:
		_, err := tx.Exec(ctx, cancelPaymentByOrderQuery, model.OrderID, model.Status)
		return err
	default:
		_, err := tx.Exec(ctx, updatePaymentQuery, model.PaymentID, model.Status)
		return err
	}
}

const (
	findPaymentQuery          = `SELECT customer_id, order_id, amount, status FROM payments WHERE payment_id = $1`
	findBalanceQuery          = `SELECT available_amount, reserved_amount FROM balances WHERE customer_id = $1`
	updateBalanceQuery        = `UPDATE balances SET available_amount = $2, reserved_amount = $3 WHERE customer_id = $1`
	insertPaymentQuery        = `INSERT INTO payments(payment_id, status, customer_id, order_id, amount) VALUES ($1, $2, $3, $4, $5)`
	updatePaymentQuery        = `UPDATE payments SET status = $2 WHERE payment_id = $1`
	cancelPaymentByOrderQuery = `UPDATE payments SET status = $2 WHERE order_id = $1`
)
