//go:build integration
// +build integration

package repository

import (
	"context"
	"os"
	"testing"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v4/log/zerologadapter"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/require"

	"github.com/moeryomenko/saga/internal/order/domain"
)

func TestIntegration_SelectQuery(t *testing.T) {
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
		},
	}

	for name, tc := range testcase {
		tc := tc
		t.Run(name, func(t *testing.T) {
			var (
				order domain.Order
				err   error
			)
			for _, event := range tc.getEvents(tc.orderID, tc.customerID) {
				order, err = PersistOrder(context.Background(), tc.orderID, event)
				require.NoError(t, err)
			}
			tc.expectedOrderState(order)
		})
	}
}
