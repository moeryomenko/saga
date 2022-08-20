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
		return service.HandleEvent(ctx, orderID, domain.CreateOrder{
			OrderID:    orderID,
			CustomerID: *createOrder.CustomerId,
		})
	}), WithResponseMapper(mapOrder), WithDefaultStatus(http.StatusCreated))
}

func (RestController) PostOrderOrderID(w http.ResponseWriter, r *http.Request, orderID openapi_types.UUID) {
	handlerDecorator(w, r, WithOperation(func(ctx context.Context) (any, error) {
		return service.HandleEvent(ctx, orderID, domain.Process{})
	}), WithResponseMapper(mapOrder), WithErrorMapper(mapDomainError))
}

func (RestController) PutOrderOrderID(w http.ResponseWriter, r *http.Request, orderID openapi_types.UUID) {
	var item Item
	handlerDecorator(w, r, WithRequestBody(&item), WithOperation(func(ctx context.Context) (any, error) {
		return service.HandleEvent(ctx, orderID, domain.AddItem{
			Item: *item.Name,
		})
	}), WithResponseMapper(mapOrder), WithErrorMapper(mapDomainError))
}

func (RestController) DeleteOrderOrderIDItem(w http.ResponseWriter, r *http.Request, orderID openapi_types.UUID, item string) {
	handlerDecorator(w, r, WithOperation(func(ctx context.Context) (any, error) {
		return service.HandleEvent(ctx, orderID, domain.RemoveItem{
			Item: item,
		})
	}), WithResponseMapper(mapOrder), WithErrorMapper(mapDomainError))
}

func mapDomainError(err error) int {
	if errors.Is(err, domain.ErrDomain) {
		return http.StatusPreconditionFailed
	}
	return http.StatusInternalServerError
}
