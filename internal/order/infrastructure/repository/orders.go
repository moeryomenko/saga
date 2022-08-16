package repository

import (
	"context"

	"github.com/gofrs/uuid/v3"
	"github.com/jackc/pgtype"
	"github.com/jackc/pgx/v4"

	"github.com/moeryomenko/saga/internal/order/domain"
	"github.com/moeryomenko/saga/pkg/errors"
)

func PersistOrder(ctx context.Context, orderID uuid.UUID, event domain.Event) (domain.Order, error) {
	var order domain.Order
	err := pool.BeginTxFunc(ctx, pgx.TxOptions{IsoLevel: pgx.ReadCommitted}, func(tx pgx.Tx) (err error) {
		order, err = findOrderByID(ctx, tx, orderID)
		switch err {
		case nil, pgx.ErrNoRows:
			// it's ok for CreateOrder event.
		default:
			return errors.MarkAndWrapError(err, ErrInfrastructure, `couldn't find order`)
		}

		order, err = domain.Apply(order, event)
		if err != nil {
			return errors.MarkAndWrapError(err, domain.ErrDomain, `couldn't apply event`)
		}

		return saveOrder(ctx, tx, order)
	})
	return order, err
}

func findOrderByID(ctx context.Context, tx pgx.Tx, orderID uuid.UUID) (domain.Order, error) {
	order := &Order{OrderID: pgtype.UUID{Bytes: orderID}}
	err := tx.QueryRow(ctx, findOrderQuery, orderID.String()).Scan(
		&order.CustomerID,
		&order.Items,
		&order.Price,
		&order.PaymentID,
		&order.Kind,
	)
	if err != nil {
		return nil, err
	}
	return mapToDomain(order), nil
}

func saveOrder(ctx context.Context, tx pgx.Tx, order domain.Order) error {
	model, err := mapToModel(order)
	if err != nil {
		return err
	}
	var query string
	switch model.Kind {
	case empty:
		query = insertOrderQuery
	default:
		query = updateOrderQuery
	}
	_, err = tx.Exec(ctx, query, model.OrderID, model.CustomerID, model.Items, model.Price, model.PaymentID, model.Kind)
	return errors.MarkAndWrapError(err, ErrInfrastructure, `couldn't save order`)
}

const (
	findOrderQuery   = `SELECT customer_id, items, price, payment_id, kind FROM orders WHERE order_id = $1 FOR UPDATE`
	insertOrderQuery = `INSERT INTO orders(order_id, customer_id, items, price, payment_id, kind) VALUES ($1, $2, $3, $4, $5, $6)`
	updateOrderQuery = `UPDATE orders SET customer_id = $2, items = $3, price = $4, payment_id = $5, kind = $6 WHERE order_id = $1`
)
