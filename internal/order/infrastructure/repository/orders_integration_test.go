package repository

import (
	"context"
	"os"
	"testing"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/log/zerologadapter"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/rs/zerolog"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/require"

	"github.com/moeryomenko/saga/internal/order/domain"
	"github.com/moeryomenko/saga/schema"
)

func TestIntegration_Repository(t *testing.T) {
	config, err := pgxpool.ParseConfig(`user=test password=pass host=localhost port=5432 dbname=orders pool_max_conns=1`)
	require.NoError(t, err)

	zlog := zerolog.New(zerolog.ConsoleWriter{Out: os.Stderr})
	config.ConnConfig.Logger = zerologadapter.NewLogger(zlog)

	pool, err = pgxpool.ConnectConfig(context.Background(), config)
	require.NoError(t, err)
	defer func() {
		pool.Close()
	}()

	testcase := map[string]struct {
		orderID, customerID uuid.UUID
		getEvents           func(orderID, customerID uuid.UUID) []domain.Event
		expectedOrderState  func(order domain.Order)
		expectedEvent       func(orderID, customerID uuid.UUID) []schema.OrderEvent
	}{
		`success order completion`: {
			orderID:    genUUID(t),
			customerID: genUUID(t),
			getEvents: func(orderID, customerID uuid.UUID) []domain.Event {
				return []domain.Event{
					domain.CreateOrder{OrderID: orderID, CustomerID: customerID},
					domain.AddItem{Item: `test`},
					domain.Process{},
					domain.ConfirmStock{},
					domain.ConfirmPayment{},
				}
			},
			expectedOrderState: func(order domain.Order) {
				if _, ok := order.(domain.CompletedOrder); !ok {
					require.FailNow(t, `expected order completed`)
				}
			},
			expectedEvent: func(orderID, customerID uuid.UUID) []schema.OrderEvent {
				return []schema.OrderEvent{
					{
						Event:      schema.Event{Type: schema.NewOrder},
						OrderID:    orderID,
						CustomerID: customerID,
						Items:      `test`,
						Price:      decimal.NewFromFloat32(9.99),
					},
					{
						Event:      schema.Event{Type: schema.CompleteOrder},
						OrderID:    orderID,
						CustomerID: customerID,
						Items:      `test`,
						Price:      decimal.NewFromFloat32(9.99),
					},
				}
			},
		},
		`success order completion (with different order confirmation)`: {
			orderID:    genUUID(t),
			customerID: genUUID(t),
			getEvents: func(orderID, customerID uuid.UUID) []domain.Event {
				return []domain.Event{
					domain.CreateOrder{OrderID: orderID, CustomerID: customerID},
					domain.AddItem{Item: `test`},
					domain.Process{},
					domain.ConfirmPayment{},
					domain.ConfirmStock{},
				}
			},
			expectedOrderState: func(order domain.Order) {
				if _, ok := order.(domain.CompletedOrder); !ok {
					require.FailNow(t, `expected order completed`)
				}
			},
			expectedEvent: func(orderID, customerID uuid.UUID) []schema.OrderEvent {
				return []schema.OrderEvent{
					{
						Event:      schema.Event{Type: schema.NewOrder},
						OrderID:    orderID,
						CustomerID: customerID,
						Items:      `test`,
						Price:      decimal.NewFromFloat32(9.99),
					},
					{
						Event:      schema.Event{Type: schema.CompleteOrder},
						OrderID:    orderID,
						CustomerID: customerID,
						Items:      `test`,
						Price:      decimal.NewFromFloat32(9.99),
					},
				}
			},
		},
		`success order completion (with two items)`: {
			orderID:    genUUID(t),
			customerID: genUUID(t),
			getEvents: func(orderID, customerID uuid.UUID) []domain.Event {
				return []domain.Event{
					domain.CreateOrder{OrderID: orderID, CustomerID: customerID},
					domain.AddItem{Item: `test`},
					domain.AddItem{Item: `test1`},
					domain.Process{},
					domain.ConfirmPayment{},
					domain.ConfirmStock{},
				}
			},
			expectedOrderState: func(order domain.Order) {
				if _, ok := order.(domain.CompletedOrder); !ok {
					require.FailNow(t, `expected order completed`)
				}
			},
			expectedEvent: func(orderID, customerID uuid.UUID) []schema.OrderEvent {
				return []schema.OrderEvent{
					{
						Event:      schema.Event{Type: schema.NewOrder},
						OrderID:    orderID,
						CustomerID: customerID,
						Items:      `test,test1`,
						Price:      decimal.NewFromFloat32(19.98),
					},
					{
						Event:      schema.Event{Type: schema.CompleteOrder},
						OrderID:    orderID,
						CustomerID: customerID,
						Items:      `test,test1`,
						Price:      decimal.NewFromFloat32(19.98),
					},
				}
			},
		},
		`success order completion (with one item, but with removing)`: {
			orderID:    genUUID(t),
			customerID: genUUID(t),
			getEvents: func(orderID, customerID uuid.UUID) []domain.Event {
				return []domain.Event{
					domain.CreateOrder{OrderID: orderID, CustomerID: customerID},
					domain.AddItem{Item: `test`},
					domain.AddItem{Item: `test1`},
					domain.RemoveItem{Item: `test1`},
					domain.Process{},
					domain.ConfirmPayment{},
					domain.ConfirmStock{},
				}
			},
			expectedOrderState: func(order domain.Order) {
				if _, ok := order.(domain.CompletedOrder); !ok {
					require.FailNow(t, `expected order completed`)
				}
			},
			expectedEvent: func(orderID, customerID uuid.UUID) []schema.OrderEvent {
				return []schema.OrderEvent{
					{
						Event:      schema.Event{Type: schema.NewOrder},
						OrderID:    orderID,
						CustomerID: customerID,
						Items:      `test`,
						Price:      decimal.NewFromFloat32(9.99),
					},
					{
						Event:      schema.Event{Type: schema.CompleteOrder},
						OrderID:    orderID,
						CustomerID: customerID,
						Items:      `test`,
						Price:      decimal.NewFromFloat32(9.99),
					},
				}
			},
		},
		`remove all items from order`: {
			orderID:    genUUID(t),
			customerID: genUUID(t),
			getEvents: func(orderID, customerID uuid.UUID) []domain.Event {
				return []domain.Event{
					domain.CreateOrder{OrderID: orderID, CustomerID: customerID},
					domain.AddItem{Item: `test`},
					domain.AddItem{Item: `test1`},
					domain.RemoveItem{Item: `test1`},
					domain.RemoveItem{Item: `test`},
				}
			},
			expectedOrderState: func(order domain.Order) {
				if _, ok := order.(domain.EmptyOrder); !ok {
					require.FailNow(t, `expected order completed`)
				}
			},
		},
		`cancel order by stock`: {
			orderID:    genUUID(t),
			customerID: genUUID(t),
			getEvents: func(orderID, customerID uuid.UUID) []domain.Event {
				return []domain.Event{
					domain.CreateOrder{OrderID: orderID, CustomerID: customerID},
					domain.AddItem{Item: `test`},
					domain.Process{},
					domain.ConfirmPayment{},
					domain.RejectStock{},
				}
			},
			expectedOrderState: func(order domain.Order) {
				if _, ok := order.(domain.CanceledOrder); !ok {
					require.FailNow(t, `expected order completed`)
				}
			},
			expectedEvent: func(orderID, customerID uuid.UUID) []schema.OrderEvent {
				return []schema.OrderEvent{
					{
						Event:      schema.Event{Type: schema.NewOrder},
						OrderID:    orderID,
						CustomerID: customerID,
						Items:      `test`,
						Price:      decimal.NewFromFloat32(9.99),
					},
					{
						Event:      schema.Event{Type: schema.CancelOrder},
						OrderID:    orderID,
						CustomerID: customerID,
						Items:      `test`,
						Price:      decimal.NewFromFloat32(9.99),
					},
				}
			},
		},
		`cancel order by stock (with different order events)`: {
			orderID:    genUUID(t),
			customerID: genUUID(t),
			getEvents: func(orderID, customerID uuid.UUID) []domain.Event {
				return []domain.Event{
					domain.CreateOrder{OrderID: orderID, CustomerID: customerID},
					domain.AddItem{Item: `test`},
					domain.Process{},
					domain.RejectStock{},
					domain.ConfirmPayment{},
				}
			},
			expectedOrderState: func(order domain.Order) {
				if _, ok := order.(domain.CanceledOrder); !ok {
					require.FailNow(t, `expected order completed`)
				}
			},
			expectedEvent: func(orderID, customerID uuid.UUID) []schema.OrderEvent {
				return []schema.OrderEvent{
					{
						Event:      schema.Event{Type: schema.NewOrder},
						OrderID:    orderID,
						CustomerID: customerID,
						Items:      `test`,
						Price:      decimal.NewFromFloat32(9.99),
					},
					{
						Event:      schema.Event{Type: schema.CancelOrder},
						OrderID:    orderID,
						CustomerID: customerID,
						Items:      `test`,
						Price:      decimal.NewFromFloat32(9.99),
					},
				}
			},
		},
		`cancel order by payment`: {
			orderID:    genUUID(t),
			customerID: genUUID(t),
			getEvents: func(orderID, customerID uuid.UUID) []domain.Event {
				return []domain.Event{
					domain.CreateOrder{OrderID: orderID, CustomerID: customerID},
					domain.AddItem{Item: `test`},
					domain.Process{},
					domain.ConfirmStock{},
					domain.RejectPayment{},
				}
			},
			expectedOrderState: func(order domain.Order) {
				if _, ok := order.(domain.CanceledOrder); !ok {
					require.FailNow(t, `expected order completed`)
				}
			},
			expectedEvent: func(orderID, customerID uuid.UUID) []schema.OrderEvent {
				return []schema.OrderEvent{
					{
						Event:      schema.Event{Type: schema.NewOrder},
						OrderID:    orderID,
						CustomerID: customerID,
						Items:      `test`,
						Price:      decimal.NewFromFloat32(9.99),
					},
					{
						Event:      schema.Event{Type: schema.CancelOrder},
						OrderID:    orderID,
						CustomerID: customerID,
						Items:      `test`,
						Price:      decimal.NewFromFloat32(9.99),
					},
				}
			},
		},
		`cancel order by payment (with different order events)`: {
			orderID:    genUUID(t),
			customerID: genUUID(t),
			getEvents: func(orderID, customerID uuid.UUID) []domain.Event {
				return []domain.Event{
					domain.CreateOrder{OrderID: orderID, CustomerID: customerID},
					domain.AddItem{Item: `test`},
					domain.Process{},
					domain.RejectPayment{},
					domain.ConfirmStock{},
				}
			},
			expectedOrderState: func(order domain.Order) {
				if _, ok := order.(domain.CanceledOrder); !ok {
					require.FailNow(t, `expected order completed`)
				}
			},
			expectedEvent: func(orderID, customerID uuid.UUID) []schema.OrderEvent {
				return []schema.OrderEvent{
					{
						Event:      schema.Event{Type: schema.NewOrder},
						OrderID:    orderID,
						CustomerID: customerID,
						Items:      `test`,
						Price:      decimal.NewFromFloat32(9.99),
					},
					{
						Event:      schema.Event{Type: schema.CancelOrder},
						OrderID:    orderID,
						CustomerID: customerID,
						Items:      `test`,
						Price:      decimal.NewFromFloat32(9.99),
					},
				}
			},
		},
	}

	for name, tc := range testcase {
		tc := tc
		t.Run(name, func(t *testing.T) {
			ctx := context.Background()
			pool.BeginTxFunc(ctx, pgx.TxOptions{IsoLevel: pgx.ReadCommitted}, func(tx pgx.Tx) (err error) {
				_, err = tx.Exec(ctx, `TRUNCATE event_log`)
				require.NoError(t, err)
				_, err = tx.Exec(ctx, `UPDATE event_offset SET offset_acked = 0`)
				require.NoError(t, err)
				return nil
			})

			var (
				order domain.Order
				err   error
			)
			for _, event := range tc.getEvents(tc.orderID, tc.customerID) {
				order, err = PersistOrder(context.Background(), tc.orderID, event)
				require.NoError(t, err)
			}
			tc.expectedOrderState(order)

			if tc.expectedEvent != nil {
				for _, expectedEvent := range tc.expectedEvent(tc.orderID, tc.customerID) {
					id, event, err := GetEvent(ctx)
					require.NoError(t, err)
					require.Equal(t, expectedEvent, event)
					err = Ack(ctx, id)
					require.NoError(t, err)
				}
			}
		})
	}
}
