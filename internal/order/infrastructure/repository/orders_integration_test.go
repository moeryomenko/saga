//go:build integration
// +build integration

package repository

import (
	"context"
	"testing"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/stretchr/testify/require"

	"github.com/moeryomenko/saga/internal/order/domain"
)

func TestIntegration_SelectQuery(t *testing.T) {
	config, err := pgxpool.ParseConfig(`user=test password=pass host=localhost port=5432 dbname=orders pool_max_conns=1`)
	require.NoError(t, err)

	pool, err = pgxpool.ConnectConfig(context.Background(), config)
	require.NoError(t, err)
	defer func() {
		pool.Close()
	}()

	ctx := context.Background()

	orderID := genUUID(t)
	customerID := genUUID(t)
	var event domain.Event
	event = domain.CreateOrder{OrderID: orderID, CustomerID: customerID}

	_, err = PersistOrder(ctx, orderID, event)
	require.NoError(t, err)

	event = domain.AddItem{Item: `test`}

	_, err = PersistOrder(ctx, orderID, event)
	require.NoError(t, err)

	var expectedOrder domain.Order = domain.ActiveOrder{
		EmptyOrder: domain.EmptyOrder{
			ID:         orderID,
			CustomerID: customerID,
		},
		Items: []string{`test`},
	}

	var order domain.Order
	pool.BeginTxFunc(ctx, pgx.TxOptions{IsoLevel: pgx.ReadCommitted}, func(tx pgx.Tx) (err error) {
		order, err = findOrderByID(ctx, tx, orderID)
		return err
	})
	require.NoError(t, err)
	require.Equal(t, expectedOrder, order)
}
