package domain

import "errors"

var (
	ErrDomain = errors.New(`domain`)

	ErrInsufficientFunds = errors.New(`insufficient funds to pay`)
	ErrCanceledPayment   = errors.New(`compelete canceled payment`)
	ErrFailedPayment     = errors.New(`cancel failed payment`)
)
