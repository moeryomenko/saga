package repository

import "errors"

var (
	ErrInfrastructure = errors.New(`infrastructure`)

	ErrNoEvents = errors.New(`no new event into log`)
)
