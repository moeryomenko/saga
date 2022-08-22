package eventhandler

import (
	"testing"

	"github.com/google/uuid"
	"github.com/moeryomenko/saga/internal/order/domain"
	"github.com/moeryomenko/saga/schema"
	"github.com/stretchr/testify/require"
)

func Test_mapToDomainEvent(t *testing.T) {
	testcases := map[string]struct {
		orderID, paymentID  uuid.UUID
		getEvent            func(orderID, paymentID uuid.UUID) map[string]any
		expectedDomainEvent func(orderID, paymentID uuid.UUID) domain.Event
	}{
		`payment confirmed`: {
			orderID: uuid.New(),
			getEvent: func(orderID, paymentID uuid.UUID) map[string]any {
				event := schema.PaymentsEvent{OrderID: orderID, PaymentsID: paymentID}
				event.SetType(schema.PaymentsConfirmed)
				return mapToEvent(event.Map())
			},
			expectedDomainEvent: func(orderID, paymentID uuid.UUID) domain.Event {
				return domain.ConfirmPayment{
					PaymentID: paymentID,
				}
			},
		},
		`payment failed`: {
			orderID: uuid.New(),
			getEvent: func(orderID, paymentID uuid.UUID) map[string]any {
				event := schema.PaymentsEvent{OrderID: orderID, PaymentsID: paymentID}
				event.SetType(schema.PaymentsFailed)
				return mapToEvent(event.Map())
			},
			expectedDomainEvent: func(orderID, paymentID uuid.UUID) domain.Event {
				return domain.RejectPayment{}
			},
		},
		`stock confirmed`: {
			orderID: uuid.New(),
			getEvent: func(orderID, paymentID uuid.UUID) map[string]any {
				event := schema.StockEvent{OrderID: orderID}
				event.SetType(schema.StockConfirmed)
				return mapToEvent(event.Map())
			},
			expectedDomainEvent: func(orderID, paymentID uuid.UUID) domain.Event {
				return domain.ConfirmStock{}
			},
		},
		`stock failed`: {
			orderID: uuid.New(),
			getEvent: func(orderID, paymentID uuid.UUID) map[string]any {
				event := schema.StockEvent{OrderID: orderID}
				event.SetType(schema.StockFailed)
				return mapToEvent(event.Map())
			},
			expectedDomainEvent: func(orderID, paymentID uuid.UUID) domain.Event {
				return domain.RejectStock{}
			},
		},
	}

	for name, tc := range testcases {
		tc := tc
		t.Run(name, func(t *testing.T) {
			orderID, event, err := mapToDomainEvent(tc.getEvent(tc.orderID, tc.paymentID))
			require.NoError(t, err)
			require.Equal(t, tc.orderID, orderID)
			require.Equal(t, tc.expectedDomainEvent(tc.orderID, tc.paymentID), event)
		})
	}
}

func mapToEvent(event map[string]string) map[string]any {
	mappedEvent := make(map[string]any)
	for key, value := range event {
		mappedEvent[key] = value
	}
	return mappedEvent
}
