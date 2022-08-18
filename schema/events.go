package schema

import (
	"encoding/json"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

func GetEventType(data map[string]any) EventType {
	kind := data[`type`].(EventType)
	return kind
}

type Event struct {
	Type EventType `json:"type"`
}

type EventType string

const (
	NewOrder          EventType = `new_order`
	CancelOrder       EventType = `cancale_order`
	CompleteOrder     EventType = `complete_order`
	PaymentsConfirmed EventType = `payments_confirmed`
	PaymentsFailed    EventType = `paymants_failed`
	StockConfirmed    EventType = `stock_confirmed`
	StockFailed       EventType = `stock_failed`
)

type OrderEvent struct {
	Event
	OrderID    uuid.UUID       `json:"order_id"`
	CustomerID uuid.UUID       `json:"customer_id"`
	Price      decimal.Decimal `json:"price"`
	PaymentID  uuid.UUID       `json:"payment_id,omitempty"`
	Items      string          `json:"items"`
}

func (e OrderEvent) Map() map[string]string {
	var m map[string]string
	in, _ := json.Marshal(e)
	_ = json.Unmarshal(in, &m)
	return m
}

func ToOrderEvent(values map[string]any) (OrderEvent, error) {
	b, err := json.Marshal(values)
	if err != nil {
		return OrderEvent{}, nil
	}
	var o OrderEvent
	err = json.Unmarshal(b, &o)
	return o, err
}

type RollbackEvent struct {
	Event
	OrderID uuid.UUID
}

func (e RollbackEvent) Map() map[string]string {
	var m map[string]string
	in, _ := json.Marshal(e)
	_ = json.Unmarshal(in, &m)
	return m
}

func ToRolbackEvent(values map[string]any) (OrderEvent, error) {
	b, err := json.Marshal(values)
	if err != nil {
		return OrderEvent{}, nil
	}
	var o OrderEvent
	err = json.Unmarshal(b, &o)
	return o, err
}

type PaymentsEvent struct {
	Event
	OrderID    uuid.UUID `json:"order_id"`
	PaymentsID uuid.UUID `json:"payments_id"`
}

func (e PaymentsEvent) Map() map[string]string {
	var m map[string]string
	in, _ := json.Marshal(e)
	_ = json.Unmarshal(in, &m)
	return m
}

func ToPaymentsEvent(values map[string]any) (PaymentsEvent, error) {
	b, err := json.Marshal(values)
	if err != nil {
		return PaymentsEvent{}, nil
	}
	var o PaymentsEvent
	err = json.Unmarshal(b, &o)
	return o, err
}

func (e *PaymentsEvent) SetType(kind EventType) {
	e.Type = kind
}

type StockEvent struct {
	Event
	OrderID uuid.UUID `json:"order_id"`
}

func (e StockEvent) Map() map[string]string {
	var m map[string]string
	in, _ := json.Marshal(e)
	_ = json.Unmarshal(in, &m)
	return m
}

func ToStockEvent(values map[string]any) (StockEvent, error) {
	b, err := json.Marshal(values)
	if err != nil {
		return StockEvent{}, nil
	}
	var o StockEvent
	err = json.Unmarshal(b, &o)
	return o, err
}

func (e *StockEvent) SetType(kind EventType) {
	e.Type = kind
}
