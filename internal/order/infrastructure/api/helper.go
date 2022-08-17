package api

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
)

type HandlerDecorator struct {
	errMapper   func(error) int
	requestBody Validated
	operation   func(context.Context) (any, int, error)
	mutation    bool
}

type Option func(*HandlerDecorator)

func WithErrorMapper(mapper func(error) int) Option {
	return func(hd *HandlerDecorator) {
		hd.errMapper = mapper
	}
}

func WithRequestBody(request Validated) Option {
	return func(hd *HandlerDecorator) {
		hd.requestBody = request
		hd.mutation = true
	}
}

func WithOperation(op func(context.Context) (any, int, error)) Option {
	return func(hd *HandlerDecorator) {
		hd.operation = op
	}
}

func handlerDecorator(w http.ResponseWriter, r *http.Request, opts ...Option) {
	decorator := &HandlerDecorator{}
	for _, opt := range opts {
		opt(decorator)
	}

	ctx := r.Context()

	if decorator.mutation {
		defer func() { _ = r.Body.Close() }()

		var body []byte
		body, err := io.ReadAll(r.Body)
		if err != nil {
			apiError(ctx, w, err.Error(), http.StatusBadRequest)
			return
		}

		err = json.Unmarshal(body, decorator.requestBody)
		if err != nil {
			apiError(ctx, w, err.Error(), http.StatusBadRequest)
			return
		}

		err = decorator.requestBody.Validate()
		if err != nil {
			apiError(ctx, w, err.Error(), http.StatusBadRequest)
			return
		}
	}

	resp, status, err := decorator.operation(ctx)
	switch err {
	case nil:
		apiSuccess(ctx, w, status, resp)
	default:
		status := http.StatusInternalServerError
		if decorator.errMapper != nil {
			status = decorator.errMapper(err)
		}
		apiError(ctx, w, err.Error(), status)
	}
}

func apiError(ctx context.Context, w http.ResponseWriter, err string, status int) {
	body, _ := json.Marshal(Error{
		Errors: &[]string{err},
	})
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(status)
	_, _ = w.Write(body)
}

func apiSuccess(ctx context.Context, w http.ResponseWriter, status int, resp any) {
	body, _ := json.Marshal(resp)
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(status)
	_, _ = w.Write(body)
}
