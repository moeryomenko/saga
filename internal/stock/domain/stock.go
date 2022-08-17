package domain

import (
	"errors"

	"github.com/google/uuid"
)

type Event any

type StockOrder struct {
	OrderID uuid.UUID
	Items   []string
}

type CancelStock struct {
	OrderID uuid.UUID
}

func Apply(event Event) (Stock, error) {
	switch event := event.(type) {
	case StockOrder:
		return CreateStock(event.OrderID, event.Items), nil
	case CancelStock:
		return CanceledStock(event), nil
	default:
		panic(`bug: invalid domain event`)
	}
}

var ErrCantStock = errors.New(`cant stock order`)

type Stock interface {
	GetOrderID() uuid.UUID
}

type CanceledStock struct {
	OrderID uuid.UUID
}

func (s CanceledStock) GetOrderID() uuid.UUID {
	return s.OrderID
}

type RejectedStock struct {
	OrderID uuid.UUID
}

func (s RejectedStock) GetOrderID() uuid.UUID {
	return s.OrderID
}

type ActiveStock struct {
	ID      uuid.UUID
	OrderID uuid.UUID
	Items   []string
}

func (s ActiveStock) GetOrderID() uuid.UUID {
	return s.OrderID
}

func CreateStock(orderID uuid.UUID, items []string) Stock {
	// dummy stock domain logic.
	if len(items) > 10 {
		return RejectedStock{OrderID: orderID}
	}

	stock := ActiveStock{
		ID:      uuid.New(),
		OrderID: orderID,
	}

	copy(stock.Items, items)

	return stock
}
