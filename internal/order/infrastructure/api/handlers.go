package api

import (
	"context"
	"errors"
	"net/http"
	"time"

	openapi_types "github.com/deepmap/oapi-codegen/pkg/types"
	"github.com/google/uuid"
	"github.com/moeryomenko/saga/internal/order/config"
	"github.com/moeryomenko/saga/internal/order/domain"
	"github.com/moeryomenko/saga/internal/order/service"
)

func New(cfg *config.Config) *http.Server {
	return &http.Server{
		ReadHeaderTimeout: 1 * time.Minute,
		Handler:           Handler(RestController{}),
		Addr:              cfg.Addr(),
	}
}

type RestController struct{}

func (RestController) PostOrder(w http.ResponseWriter, r *http.Request) {
	var createOrder CreateOrder
	handlerDecorator(w, r, WithRequestBody(&createOrder), WithOperation(func(ctx context.Context) (any, error) {
		orderID := uuid.New()

		order, err := service.HandleEvent(ctx, orderID, domain.CreateOrder{
			OrderID:    orderID,
			CustomerID: *createOrder.CustomerId,
		})
		if err != nil {
			return nil, err
		}

		return mapOrder(order), nil
	}), WithDefaultStatus(http.StatusCreated))
}

func (RestController) PostOrderOrderID(w http.ResponseWriter, r *http.Request, orderID openapi_types.UUID) {
	handlerDecorator(w, r, WithOperation(func(ctx context.Context) (any, error) {
		_, err := service.HandleEvent(ctx, orderID, domain.Process{})
		return nil, err
	}), WithErrorMapper(mapDomainError))
}

func (RestController) PutOrderOrderID(w http.ResponseWriter, r *http.Request, orderID openapi_types.UUID) {
	var item Item
	handlerDecorator(w, r, WithRequestBody(&item), WithOperation(func(ctx context.Context) (any, error) {
		_, err := service.HandleEvent(ctx, orderID, domain.AddItem{
			Item: *item.Name,
		})
		return nil, err
	}), WithDefaultStatus(http.StatusNoContent), WithErrorMapper(mapDomainError))
}

func (RestController) DeleteOrderOrderIDItem(w http.ResponseWriter, r *http.Request, orderID openapi_types.UUID, item string) {
	handlerDecorator(w, r, WithOperation(func(ctx context.Context) (any, error) {
		_, err := service.HandleEvent(ctx, orderID, domain.RemoveItem{
			Item: item,
		})
		return nil, err
	}), WithErrorMapper(mapDomainError))
}

func mapDomainError(err error) int {
	if errors.Is(err, domain.ErrDomain) {
		return http.StatusPreconditionFailed
	}
	return http.StatusInternalServerError
}
