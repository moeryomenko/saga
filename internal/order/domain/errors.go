package domain

import "errors"

var (
	ErrDomain = errors.New(`domain`)

	ErrCancelOrder     = errors.New(`cancellation of an unfinished order`)
	ErrAddItem         = errors.New(`adding a item to an order that is being processed`)
	ErrRemoveItemEmpty = errors.New(`remove item from empty order`)
	ErrRemoveItem      = errors.New(`removing a item from an order that is being processed`)
	ErrPayOrder        = errors.New(`payment for a prepared order`)
	ErrStockOrder      = errors.New(`stocking of not prepared order`)
	ErrEmptyOrder      = errors.New(`couldn't process empty order`)
)
