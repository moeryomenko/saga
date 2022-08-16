package domain

import (
	uuid "github.com/gofrs/uuid/v3"
	"github.com/shopspring/decimal"
)

type Order interface {
	GetID() uuid.UUID
}

type EmptyOrder struct {
	// ID is id of order.
	ID uuid.UUID
	// CustomerID needs for reference to the customer.
	CustomerID uuid.UUID
}

func (o EmptyOrder) GetID() uuid.UUID {
	return o.ID
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

type StockedOrder struct {
	PendingOrder
}

type CanceledOrder struct {
	PendingOrder
}

type PaidOrder struct {
	PendingOrder

	// PaymentID needs for reference to the payment.
	PaymentID uuid.UUID
}

type CompletedOrder struct {
	PaidOrder
}

// AttachPayments attachs payments info to order.
func AttachPayments(order Order, paymentID uuid.UUID) (Order, error) {
	switch order := order.(type) {
	case PendingOrder:
		return PaidOrder{
			PendingOrder: order,
			PaymentID:    paymentID,
		}, nil
	case StockedOrder:
		return CompletedOrder{
			PaidOrder: PaidOrder{
				PendingOrder: order.PendingOrder,
				PaymentID:    paymentID,
			},
		}, nil
	default:
		return nil, ErrPayOrder
	}
}

// StockOrder mark order as stocked.
func StockOrder(order Order) (Order, error) {
	switch order := order.(type) {
	case PendingOrder:
		return StockedOrder{
			PendingOrder: order,
		}, nil
	case PaidOrder:
		return CompletedOrder{
			PaidOrder: order,
		}, nil
	default:
		return nil, ErrStockOrder
	}
}

// CancelOrder cancels order.
func CancelOrder(order Order) (CanceledOrder, error) {
	switch order := order.(type) {
	case PaidOrder:
		return CanceledOrder{PendingOrder: order.PendingOrder}, nil
	case StockedOrder:
		return CanceledOrder(order), nil
	default:
		return CanceledOrder{}, ErrCancelOrder
	}
}

// AddItemToOrder adds item to order.
func AddItemToOrder(order Order, item string) (ActiveOrder, error) {
	switch order := any(order).(type) {
	case EmptyOrder:
		return ActiveOrder{
			EmptyOrder: EmptyOrder{ID: order.ID, CustomerID: order.CustomerID},
			Items:      append([]string{}, item),
		}, nil
	case ActiveOrder:
		return ActiveOrder{
			EmptyOrder: EmptyOrder{ID: order.ID, CustomerID: order.CustomerID},
			Items:      append(order.Items, item),
		}, nil
	default:
		return ActiveOrder{}, ErrAddItem
	}
}

// RemoveItemFromOrder remove given item from order.
func RemoveItemFromOrder(order Order, item string) (Order, error) {
	switch order := order.(type) {
	case EmptyOrder:
		return nil, ErrRemoveItemEmpty
	case ActiveOrder:
		items := removeItem(order.Items, item)

		if len(items) == 0 {
			return EmptyOrder{
				ID:         order.ID,
				CustomerID: order.CustomerID,
			}, nil
		}

		return ActiveOrder{
			EmptyOrder: order.EmptyOrder,
			Items:      items,
		}, nil
	default:
		return nil, ErrRemoveItem
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
