package domain

import uuid "github.com/gofrs/uuid/v3"

type Event interface{}

type CreateOrder struct {
	CustomerID uuid.UUID
}

type AddItem struct {
	Item string
}

type RemoveItem struct {
	Item string
}

type CofirmPayment struct {
	PaymentID uuid.UUID
}

type ConfirmStock struct{}

type RejectPayment struct{}

type RejectStock struct{}

func Apply(order Order, event Event) (Order, error) {
	switch event := event.(type) {
	case CreateOrder:
		orderID, _ := uuid.NewV4()
		return EmptyOrder{
			ID:         orderID,
			CustomerID: event.CustomerID,
		}, nil
	case AddItem:
		return AddItemToOrder(order, event.Item)
	case RemoveItem:
		return RemoveItemFromOrder(order, event.Item)
	case CofirmPayment:
		return AttachPayments(order, event.PaymentID)
	case ConfirmStock:
		return StockOrder(order)
	case RejectPayment, RejectStock:
		return CancelOrder(order)
	default:
		panic(`bug: invalid event type`)
	}
}
