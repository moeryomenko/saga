package repository

import (
	"errors"

	"github.com/gofrs/uuid/v3"
	"github.com/moeryomenko/saga/internal/order/domain"
	"github.com/shopspring/decimal"
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
	OrderID    uuid.UUID
	CustomerID uuid.UUID
	Items      []string
	Price      *decimal.Decimal
	PaymentID  *uuid.UUID
	Kind       string
}

func mapToDomain(o *Order) domain.Order {
	if o == nil {
		return nil
	}

	switch o.Kind {
	case empty:
		return domain.EmptyOrder{
			ID:         o.OrderID,
			CustomerID: o.CustomerID,
		}
	case active:
		return domain.ActiveOrder{
			EmptyOrder: domain.EmptyOrder{
				ID:         o.OrderID,
				CustomerID: o.CustomerID,
			},
			Items: o.Items,
		}
	case pending:
		return domain.PendingOrder{
			ActiveOrder: domain.ActiveOrder{
				EmptyOrder: domain.EmptyOrder{
					ID:         o.OrderID,
					CustomerID: o.CustomerID,
				},
				Items: o.Items,
			},
			Price: *o.Price,
		}
	case stocked:
		return domain.StockedOrder{
			PendingOrder: domain.PendingOrder{
				ActiveOrder: domain.ActiveOrder{
					EmptyOrder: domain.EmptyOrder{
						ID:         o.OrderID,
						CustomerID: o.CustomerID,
					},
					Items: o.Items,
				},
				Price: *o.Price,
			},
		}
	case paid:
		return domain.PaidOrder{
			PendingOrder: domain.PendingOrder{
				ActiveOrder: domain.ActiveOrder{
					EmptyOrder: domain.EmptyOrder{
						ID:         o.OrderID,
						CustomerID: o.CustomerID,
					},
					Items: o.Items,
				},
				Price: *o.Price,
			},
			PaymentID: *o.PaymentID,
		}
	case complited:
		return domain.CompletedOrder{
			PaidOrder: domain.PaidOrder{
				PendingOrder: domain.PendingOrder{
					ActiveOrder: domain.ActiveOrder{
						EmptyOrder: domain.EmptyOrder{
							ID:         o.OrderID,
							CustomerID: o.CustomerID,
						},
						Items: o.Items,
					},
					Price: *o.Price,
				},
				PaymentID: *o.PaymentID,
			},
		}
	case canceled:
		return domain.CanceledOrder{
			PendingOrder: domain.PendingOrder{
				ActiveOrder: domain.ActiveOrder{
					EmptyOrder: domain.EmptyOrder{
						ID:         o.OrderID,
						CustomerID: o.CustomerID,
					},
					Items: o.Items,
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
		OrderID:    o.GetID(),
		CustomerID: o.GetCustomerID(),
	}

	switch o := o.(type) {
	case domain.EmptyOrder:
	case domain.ActiveOrder:
		copy(order.Items, o.Items)
	case domain.PendingOrder:
		copy(order.Items, o.Items)
		order.Price = &o.Price
	case domain.StockedOrder:
		copy(order.Items, o.Items)
		order.Price = &o.Price
	case domain.PaidOrder:
		copy(order.Items, o.Items)
		order.Price = &o.Price
		order.PaymentID = &o.PaymentID
	case domain.CompletedOrder:
		copy(order.Items, o.Items)
		order.Price = &o.Price
		order.PaymentID = &o.PaymentID
	case domain.CanceledOrder:
		copy(order.Items, o.Items)
		order.Price = &o.Price
	}

	return order, nil
}
