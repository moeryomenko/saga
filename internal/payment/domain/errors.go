package domain

import "errors"

var (
	ErrDomain = errors.New(`domain`)

	ErrInsufficientFunds = errors.New(`insufficient funds to pay`)
	ErrCanceledPayment   = errors.New(`compelete canceled payment`)
	ErrCompletedPayment  = errors.New(`cancel completed payment`)
)
