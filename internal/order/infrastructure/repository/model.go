package repository

import (
	"encoding/json"
	"errors"

	"github.com/jackc/pgtype"
	"github.com/shopspring/decimal"

	"github.com/moeryomenko/saga/internal/order/domain"
)

const (
	empty     = `empty`
	active    = `active`
	pending   = `pending`
	stocked   = `stocked`
	paid      = `paid`
	complited = `completed`
	canceled  = `canceled`
)

type Order struct {
	OrderID    pgtype.UUID
	CustomerID pgtype.UUID
	Items      pgtype.JSONB
	Price      *decimal.Decimal
	PaymentID  pgtype.UUID
	Kind       string
}

func itemsMap(i pgtype.JSONB) []string {
	var items []string
	_ = json.Unmarshal(i.Bytes, &items)
	return items
}

func itemsToModel(items []string) pgtype.JSONB {
	b, err := json.Marshal(items)
	if err != nil {
		return pgtype.JSONB{Status: pgtype.Null}
	}
	return pgtype.JSONB{Bytes: b, Status: pgtype.Present}
}

func mapToDomain(o *Order) domain.Order {
	if o == nil {
		return nil
	}

	switch o.Kind {
	case empty:
		return domain.EmptyOrder{
			ID:         o.OrderID.Bytes,
			CustomerID: o.CustomerID.Bytes,
		}
	case active:
		return domain.ActiveOrder{
			EmptyOrder: domain.EmptyOrder{
				ID:         o.OrderID.Bytes,
				CustomerID: o.CustomerID.Bytes,
			},
			Items: itemsMap(o.Items),
		}
	case pending:
		return domain.PendingOrder{
			ActiveOrder: domain.ActiveOrder{
				EmptyOrder: domain.EmptyOrder{
					ID:         o.OrderID.Bytes,
					CustomerID: o.CustomerID.Bytes,
				},
				Items: itemsMap(o.Items),
			},
			Price: *o.Price,
		}
	case stocked:
		return domain.StockedOrder{
			PendingOrder: domain.PendingOrder{
				ActiveOrder: domain.ActiveOrder{
					EmptyOrder: domain.EmptyOrder{
						ID:         o.OrderID.Bytes,
						CustomerID: o.CustomerID.Bytes,
					},
					Items: itemsMap(o.Items),
				},
				Price: *o.Price,
			},
		}
	case paid:
		return domain.PaidOrder{
			PendingOrder: domain.PendingOrder{
				ActiveOrder: domain.ActiveOrder{
					EmptyOrder: domain.EmptyOrder{
						ID:         o.OrderID.Bytes,
						CustomerID: o.CustomerID.Bytes,
					},
					Items: itemsMap(o.Items),
				},
				Price: *o.Price,
			},
			PaymentID: o.PaymentID.Bytes,
		}
	case complited:
		return domain.CompletedOrder{
			PaidOrder: domain.PaidOrder{
				PendingOrder: domain.PendingOrder{
					ActiveOrder: domain.ActiveOrder{
						EmptyOrder: domain.EmptyOrder{
							ID:         o.OrderID.Bytes,
							CustomerID: o.CustomerID.Bytes,
						},
						Items: itemsMap(o.Items),
					},
					Price: *o.Price,
				},
				PaymentID: o.PaymentID.Bytes,
			},
		}
	case canceled:
		return domain.CanceledOrder{
			PendingOrder: domain.PendingOrder{
				ActiveOrder: domain.ActiveOrder{
					EmptyOrder: domain.EmptyOrder{
						ID:         o.OrderID.Bytes,
						CustomerID: o.CustomerID.Bytes,
					},
					Items: itemsMap(o.Items),
				},
				Price: *o.Price,
			},
		}
	}

	return nil
}

func mapToModel(o domain.Order) (*Order, error) {
	if o == nil {
		return nil, errors.New(`invalid order`)
	}

	order := &Order{
		OrderID:    pgtype.UUID{Bytes: o.GetID(), Status: pgtype.Present},
		CustomerID: pgtype.UUID{Bytes: o.GetCustomerID(), Status: pgtype.Present},
		Items:      pgtype.JSONB{Status: pgtype.Null},
		PaymentID:  pgtype.UUID{Status: pgtype.Null},
	}

	switch o := any(o).(type) {
	case domain.EmptyOrder:
		order.Kind = empty
	case domain.ActiveOrder:
		order.Items = itemsToModel(o.Items)
		order.Kind = active
	case domain.PendingOrder:
		order.Items = itemsToModel(o.Items)
		order.Price = &o.Price
		order.Kind = pending
	case domain.StockedOrder:
		order.Items = itemsToModel(o.Items)
		order.Price = &o.Price
		order.Kind = stocked
	case domain.PaidOrder:
		order.Items = itemsToModel(o.Items)
		order.Price = &o.Price
		order.PaymentID = pgtype.UUID{Bytes: o.PaymentID, Status: pgtype.Present}
		order.Kind = paid
	case domain.CompletedOrder:
		order.Items = itemsToModel(o.Items)
		order.Price = &o.Price
		order.PaymentID = pgtype.UUID{Bytes: o.PaymentID, Status: pgtype.Present}
		order.Kind = complited
	case domain.CanceledOrder:
		order.Items = itemsToModel(o.Items)
		order.Price = &o.Price
		order.Kind = canceled
	default:
		return nil, errors.New(`invalid order`)
	}

	return order, nil
}
