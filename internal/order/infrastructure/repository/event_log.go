package repository

import (
	"context"
	"encoding/json"
	"strings"

	"github.com/jackc/pgtype"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/moeryomenko/saga/internal/order/domain"
	"github.com/moeryomenko/saga/pkg/errors"
	"github.com/moeryomenko/saga/schema"
)

func mapToEvent(order domain.Order) (schema.OrderEvent, bool) {
	event := schema.OrderEvent{
		OrderID:    order.GetID(),
		CustomerID: order.GetCustomerID(),
	}

	switch order := order.(type) {
	case domain.PendingOrder:
		event.SetType(schema.NewOrder)
		event.Items = strings.Join(order.Items, `,`)
		event.Price = order.Price
	case domain.CompletedOrder:
		event.SetType(schema.CompleteOrder)
		event.Items = strings.Join(order.Items, `,`)
		event.Price = order.Price
	case domain.CanceledOrder:
		event.SetType(schema.CancelOrder)
		event.Items = strings.Join(order.Items, `,`)
		event.Price = order.Price
	default:
		return schema.OrderEvent{}, false
	}

	return event, true
}

func insertEvent(ctx context.Context, tx pgx.Tx, order domain.Order) error {
	event, ok := mapToEvent(order)
	if !ok {
		return nil
	}
	payload, err := json.Marshal(event)

	_, err = tx.Exec(ctx, insertEventToLog, pgtype.JSONB{Bytes: payload, Status: pgtype.Present}, event.Type)
	return errors.MarkAndWrapError(err, ErrInfrastructure, `couldn't insert event to log`)
}

func GetEvent(ctx context.Context) (offset int, event schema.OrderEvent, err error) {
	var (
		payload   pgtype.JSONB
		eventType schema.EventType
	)
	err = pool.AcquireFunc(ctx, func(conn *pgxpool.Conn) error {
		return conn.QueryRow(ctx, selectEventFromLog).Scan(&offset, &payload, &eventType)
	})
	if err != nil {
		return 0, schema.OrderEvent{}, errors.MarkAndWrapError(err, ErrInfrastructure, `couldn't get event`)
	}
	err = json.Unmarshal(payload.Bytes, &event)
	if err != nil {
		return 0, schema.OrderEvent{}, errors.MarkAndWrapError(err, ErrInfrastructure, `invalid event payload`)
	}
	event.SetType(eventType)
	return offset, event, nil
}

func Ack(ctx context.Context, id int) error {
	err := pool.BeginTxFunc(ctx, pgx.TxOptions{IsoLevel: pgx.ReadCommitted}, func(tx pgx.Tx) (err error) {
		_, err = tx.Exec(ctx, submitOffset, id)
		return err
	})
	if err != nil {
		return errors.MarkAndWrapError(err, ErrInfrastructure, `couldn't submit offset`)
	}
	return nil
}

const (
	insertEventToLog   = `INSERT INTO event_log(payload, event_kind) VALUES ($1, $2)`
	selectEventFromLog = `
	SELECT id, payload, event_kind
	FROM event_log
	WHERE id > (SELECT offset_acked FROM event_offset)
	ORDER BY id ASC LIMIT 1`
	submitOffset = `UPDATE event_offset SET offset_acked = $1`
)
