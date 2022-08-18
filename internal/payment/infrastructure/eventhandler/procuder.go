package eventhandler

import (
	"context"

	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
	"github.com/moeryomenko/saga/internal/payment/domain"
	"github.com/moeryomenko/saga/schema"
)

func Produce(ctx context.Context, orderID uuid.UUID, payment domain.Payment) error {
	event := schema.PaymentsEvent{
		OrderID:    orderID,
		PaymentsID: payment.GetID(),
	}
	switch payment.(type) {
	case domain.NewPayment:
		event.SetType(schema.PaymentsConfirmed)
	case domain.CanceledPayment:
		event.SetType(schema.PaymentsFailed)
	}
	_, err := client.XAdd(ctx, &redis.XAddArgs{
		Stream: ConfirmStream,
		Values: event.Map(),
	}).Result()
	return err
}
