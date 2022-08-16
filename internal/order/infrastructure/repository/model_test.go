package repository

import (
	"testing"

	"github.com/gofrs/uuid/v3"
	"github.com/moeryomenko/saga/internal/order/domain"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/require"
)

func Test_mapToModel(t *testing.T) {
	testcases := []struct {
		order    domain.Order
		expected string
	}{
		{
			expected: empty,
			order: domain.EmptyOrder{
				ID:         genUUID(t),
				CustomerID: genUUID(t),
			},
		},
		{
			expected: active,
			order: domain.ActiveOrder{
				EmptyOrder: domain.EmptyOrder{
					ID:         genUUID(t),
					CustomerID: genUUID(t),
				},
				Items: []string{`test`},
			},
		},
		{
			expected: pending,
			order: domain.PendingOrder{
				ActiveOrder: domain.ActiveOrder{
					EmptyOrder: domain.EmptyOrder{
						ID:         genUUID(t),
						CustomerID: genUUID(t),
					},
					Items: []string{`test`},
				},
				Price: decimal.NewFromFloat32(9.99),
			},
		},
		{
			expected: stocked,
			order: domain.StockedOrder{
				PendingOrder: domain.PendingOrder{
					ActiveOrder: domain.ActiveOrder{
						EmptyOrder: domain.EmptyOrder{
							ID:         genUUID(t),
							CustomerID: genUUID(t),
						},
						Items: []string{`test`},
					},
					Price: decimal.NewFromFloat32(9.99),
				},
			},
		},
		{
			expected: paid,
			order: domain.PaidOrder{
				PendingOrder: domain.PendingOrder{
					ActiveOrder: domain.ActiveOrder{
						EmptyOrder: domain.EmptyOrder{
							ID:         genUUID(t),
							CustomerID: genUUID(t),
						},
						Items: []string{`test`},
					},
					Price: decimal.NewFromFloat32(9.99),
				},
				PaymentID: genUUID(t),
			},
		},
		{
			expected: complited,
			order: domain.CompletedOrder{
				PaidOrder: domain.PaidOrder{
					PendingOrder: domain.PendingOrder{
						ActiveOrder: domain.ActiveOrder{
							EmptyOrder: domain.EmptyOrder{
								ID:         genUUID(t),
								CustomerID: genUUID(t),
							},
							Items: []string{`test`},
						},
						Price: decimal.NewFromFloat32(9.99),
					},
					PaymentID: genUUID(t),
				},
			},
		},
		{
			expected: canceled,
			order: domain.CanceledOrder{
				PendingOrder: domain.PendingOrder{
					ActiveOrder: domain.ActiveOrder{
						EmptyOrder: domain.EmptyOrder{
							ID:         genUUID(t),
							CustomerID: genUUID(t),
						},
						Items: []string{`test`},
					},
					Price: decimal.NewFromFloat32(9.99),
				},
			},
		},
	}

	for _, tc := range testcases {
		tc := tc
		t.Run(tc.expected, func(t *testing.T) {
			order, err := mapToModel(tc.order)
			require.NoError(t, err)
			require.Equal(t, tc.expected, order.Kind)
			domainOrder := mapToDomain(order)
			require.Equal(t, tc.order, domainOrder)
		})
	}
}

func genUUID(t *testing.T) uuid.UUID {
	val, err := uuid.NewV4()
	require.NoError(t, err)
	return val
}
