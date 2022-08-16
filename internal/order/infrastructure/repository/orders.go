package repository

import (
	"context"

	"github.com/gofrs/uuid/v3"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"

	"github.com/moeryomenko/saga/internal/order/domain"
	"github.com/moeryomenko/saga/pkg/errors"
)

type Apply func(domain.Order, domain.Event) (domain.Order, error)

func PersistOrder(ctx context.Context, orderID uuid.UUID, event domain.Event, apply Apply) error {
	return pool.AcquireFunc(ctx, func(c *pgxpool.Conn) error {
		return c.BeginTxFunc(ctx, pgx.TxOptions{IsoLevel: pgx.ReadCommitted}, func(tx pgx.Tx) error {
			order, err := findOrderByID(ctx, tx, orderID)
			switch err {
			case pgx.ErrNoRows:
				// it's ok for CreateOrder event, and throw error for other events.
				order = nil
			default:
				return err
			}

			order, err = apply(order, event)
			if err != nil {
				return errors.MarkAndWrapError(err, domain.ErrDomain, `couldn't apply event`)
			}

			return saveOrder(ctx, tx, order)
		})
	})
}

func findOrderByID(ctx context.Context, tx pgx.Tx, orderID uuid.UUID) (domain.Order, error) {
	order := &Order{OrderID: orderID}
	err := tx.QueryRow(ctx, findOrderQuery).Scan(
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
	return nil
}

const (
	findOrderQuery = `SELECT customer_id, items, price, payment_id, kind FROM orders WHERE order_id`
)
