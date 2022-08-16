package domain

import (
	uuid "github.com/gofrs/uuid/v3"
	"github.com/shopspring/decimal"
)

type EmptyOrder struct {
	// ID is id of order.
	ID uuid.UUID
	// CustomerID needs for reference to the customer.
	CustomerID uuid.UUID
}

type ActiveOrder struct {
	EmptyOrder

	// Items is list of items into order.
	Items []string
}

type PendingOrder struct {
	ActiveOrder

	// Price is sum of items price.
	Price decimal.Decimal
}

type PaidOrder struct {
	PendingOrder

	// PaymentID needs for reference to the payment.
	PaymentID uuid.UUID
}

// AddItemToOrder adds item to order.
func AddItemToOrder[Order EmptyOrder | ActiveOrder](order Order, item string) ActiveOrder {
	switch order := any(order).(type) {
	case EmptyOrder:
		return ActiveOrder{
			EmptyOrder: EmptyOrder{ID: order.ID, CustomerID: order.CustomerID},
			Items:      append([]string{}, item),
		}
	case ActiveOrder:
		return ActiveOrder{
			EmptyOrder: EmptyOrder{ID: order.ID, CustomerID: order.CustomerID},
			Items:      append(order.Items, item),
		}
	default:
		panic(`invalid order type`)
	}
}

// RemoveItemFromOrder remove given item from order.
func RemoveItemFromOrder(order ActiveOrder, item string) (*EmptyOrder, *ActiveOrder) {
	items := removeItem(order.Items, item)

	if len(items) == 0 {
		return &EmptyOrder{
			ID:         order.ID,
			CustomerID: order.CustomerID,
		}, nil
	}

	return nil, &ActiveOrder{
		EmptyOrder: order.EmptyOrder,
		Items:      items,
	}
}

// CalculatePrice calculate price and close order to changes.
func CalculatePrice(order ActiveOrder) PendingOrder {
	return PendingOrder{
		ActiveOrder: ActiveOrder{
			EmptyOrder: order.EmptyOrder,
			Items:      order.Items,
		},
		Price: Price(order.Items),
	}
}

// Price retuns sum price of items.
func Price(items []string) decimal.Decimal {
	return decimal.NewFromFloat32(float32(len(items)) * 9.99) // dummy price.
}

func removeItem(items []string, item string) []string {
	for i, it := range items {
		if it == item {
			return append(items[:i], items[i+1:]...)
		}
	}
	return items
}
