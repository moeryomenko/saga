package repository

import (
	"context"
	"encoding/json"

	"github.com/jackc/pgtype"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/moeryomenko/saga/internal/payment/domain"
	"github.com/moeryomenko/saga/pkg/errors"
	"github.com/moeryomenko/saga/schema"
)

func GetEvent(ctx context.Context) (offset int, event schema.PaymentsEvent, err error) {
	var payload pgtype.JSONB
	err = pool.AcquireFunc(ctx, func(conn *pgxpool.Conn) error {
		return conn.QueryRow(ctx, selectEventFromLog).Scan(&offset, &payload)
	})
	switch err {
	case nil:
	case pgx.ErrNoRows:
		return 0, schema.PaymentsEvent{}, ErrNoEvents
	default:
		return 0, schema.PaymentsEvent{}, errors.MarkAndWrapError(err, ErrInfrastructure, `couldn't get event`)
	}
	err = json.Unmarshal(payload.Bytes, &event)
	if err != nil {
		return 0, schema.PaymentsEvent{}, errors.MarkAndWrapError(err, ErrInfrastructure, `invalid event payload`)
	}
	return offset, event, nil
}

func Ack(ctx context.Context, offset int) error {
	err := pool.BeginTxFunc(ctx, pgx.TxOptions{IsoLevel: pgx.ReadCommitted}, func(tx pgx.Tx) (err error) {
		_, err = tx.Exec(ctx, submitOffset, offset)
		return err
	})
	if err != nil {
		return errors.MarkAndWrapError(err, ErrInfrastructure, `couldn't submit offset`)
	}
	return nil
}

func insertEvent(ctx context.Context, tx pgx.Tx, payment domain.Payment) error {
	event, ok := mapToEvent(payment)
	if !ok {
		return nil
	}
	payload, err := json.Marshal(event)
	if err != nil {
		return errors.MarkAndWrapError(err, ErrInfrastructure, `invalid event type`)
	}

	_, err = tx.Exec(ctx, insertEventToLog, pgtype.JSONB{Bytes: payload, Status: pgtype.Present})
	if err != nil {
		return errors.MarkAndWrapError(err, ErrInfrastructure, `coudln't insert payment event to log`)
	}
	return nil
}

func mapToEvent(payment domain.Payment) (schema.PaymentsEvent, bool) {
	event := schema.PaymentsEvent{}
	switch payment := payment.(type) {
	case domain.NewPayment:
		event.OrderID = payment.OrderID
		event.PaymentsID = payment.ID
		event.SetType(schema.PaymentsConfirmed)
	case domain.FailedPayment:
		event.OrderID = payment.OrderID
		event.SetType(schema.PaymentsFailed)
	default:
		return schema.PaymentsEvent{}, false
	}
	return event, true
}

const (
	insertEventToLog   = `INSERT INTO event_log(payload) VALUES ($1)`
	selectEventFromLog = `
	SELECT id, payload
	FROM event_log
	WHERE id > (SELECT offset_acked FROM event_offset)
	ORDER BY id ASC LIMIT 1`
	submitOffset = `UPDATE event_offset SET offset_acked = $1`
)
