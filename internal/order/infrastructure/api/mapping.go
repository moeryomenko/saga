package api

import "github.com/moeryomenko/saga/internal/order/domain"

func mapOrder(order domain.Order) Order {
	switch order := order.(type) {
	case domain.EmptyOrder:
		return Order{
			Id:         &order.ID,
			CustomerId: &order.CustomerID,
		}
	case domain.ActiveOrder:
		return Order{
			Id:         &order.ID,
			CustomerId: &order.CustomerID,
			Items:      &order.Items,
		}
	case domain.PendingOrder:
		return Order{
			Id:         &order.ID,
			CustomerId: &order.CustomerID,
			Items:      &order.Items,
		}
	case domain.PaidOrder:
		return Order{
			Id:         &order.ID,
			CustomerId: &order.CustomerID,
			Items:      &order.Items,
			PaymentId:  &order.PaymentID,
		}
	case domain.StockedOrder:
		return Order{
			Id:         &order.ID,
			CustomerId: &order.CustomerID,
			Items:      &order.Items,
		}
	case domain.CompletedOrder:
		return Order{
			Id:         &order.ID,
			CustomerId: &order.CustomerID,
			Items:      &order.Items,
			PaymentId:  &order.PaymentID,
		}
	case domain.CanceledOrder:
		return Order{
			Id:         &order.ID,
			CustomerId: &order.CustomerID,
			Items:      &order.Items,
		}
	}
	return Order{}
}
