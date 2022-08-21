package domain

import "github.com/google/uuid"

type Event interface{}

type CreateOrder struct {
	OrderID    uuid.UUID
	CustomerID uuid.UUID
}

type AddItem struct {
	Item string
}

type RemoveItem struct {
	Item string
}

type Process struct{}

type ConfirmPayment struct {
	PaymentID uuid.UUID
}

type ConfirmStock struct{}

type RejectPayment struct{}

type RejectStock struct{}

func Apply(order Order, event Event) (Order, error) {
	switch event := event.(type) {
	case CreateOrder:
		return EmptyOrder{
			ID:         event.OrderID,
			CustomerID: event.CustomerID,
		}, nil
	case AddItem:
		return AddItemToOrder(order, event.Item)
	case RemoveItem:
		return RemoveItemFromOrder(order, event.Item)
	case Process:
		return CalculatePrice(order)
	case ConfirmPayment:
		return AttachPayments(order, event.PaymentID)
	case ConfirmStock:
		return StockOrder(order)
	case RejectPayment, RejectStock:
		return CancelOrder(order)
	default:
		panic(`bug: invalid event type`)
	}
}
