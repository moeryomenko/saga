//go:build integration
// +build integration

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

	"github.com/moeryomenko/saga/internal/payment/domain"
)

func TestIntegration_SelectQuery(t *testing.T) {
	config, err := pgxpool.ParseConfig(`user=test password=pass host=localhost port=5432 dbname=payments pool_max_conns=1`)
	require.NoError(t, err)

	zlog := zerolog.New(zerolog.ConsoleWriter{Out: os.Stderr})
	config.ConnConfig.Logger = zerologadapter.NewLogger(zlog)

	pool, err = pgxpool.ConnectConfig(context.Background(), config)
	require.NoError(t, err)
	defer func() {
		pool.Close()
	}()

	ctx := context.Background()

	// prepare database data.
	customerID := uuid.New()
	available := decimal.NewFromInt32(100)
	err = pool.BeginTxFunc(ctx, pgx.TxOptions{IsoLevel: pgx.ReadCommitted}, func(tx pgx.Tx) error {
		_, err := tx.Exec(ctx, `INSERT INTO balances(customer_id, available_amount) VALUES ($1, $2)`, customerID, available)
		return err
	})
	require.NoError(t, err)

	// test.

	// create payments.
	orderID := uuid.New()
	amount := decimal.NewFromInt32(20)
	var paymentID uuid.UUID
	var expectedPayment domain.Payment = domain.NewPayment{OrderID: orderID, Amount: amount}
	payment, err := PersistTransaction(ctx, customerID, domain.Reserve{OrderID: orderID, Amount: amount})
	require.NoError(t, err)
	require.Equal(t, expectedPayment.GetAmount(), payment.GetAmount())
	switch payment := payment.(type) {
	case domain.NewPayment:
		require.Equal(t, orderID, payment.OrderID)
		paymentID = payment.ID
	default:
		require.Fail(t, `invalid payment type, expected NewPayment`)
	}
	// check balance.
	expectedAvailabe := decimal.NewFromInt32(80)
	expectedReserved := decimal.NewFromInt32(20)
	var balance domain.Balance
	pool.BeginTxFunc(ctx, pgx.TxOptions{IsoLevel: pgx.ReadCommitted}, func(tx pgx.Tx) (err error) {
		balance, err = findBalanceByCustomer(ctx, tx, customerID)
		return
	})
	require.Equal(t, expectedAvailabe.String(), balance.Amount.String())
	require.Equal(t, expectedReserved.String(), balance.Reserved.String())

	// complete payments.
	expectedReserved = decimal.Zero
	expectedPayment = domain.CompletedPayment{ID: paymentID, Amount: amount}
	payment, err = PersistTransaction(ctx, customerID, domain.Complete{PaymentID: paymentID})
	require.NoError(t, err)
	switch payment := payment.(type) {
	case domain.CompletedPayment:
		require.Equal(t, paymentID, payment.ID)
		paymentID = payment.ID
	default:
		require.Fail(t, `invalid payment type, expected CompeletedPayment`)
	}
	// check balance.
	pool.BeginTxFunc(ctx, pgx.TxOptions{IsoLevel: pgx.ReadCommitted}, func(tx pgx.Tx) (err error) {
		balance, err = findBalanceByCustomer(ctx, tx, customerID)
		return
	})
	require.Equal(t, expectedAvailabe.String(), balance.Amount.String())
	require.Equal(t, expectedReserved.String(), balance.Reserved.String())
}
