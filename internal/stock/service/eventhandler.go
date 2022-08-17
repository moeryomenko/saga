package service

import "github.com/moeryomenko/saga/internal/stock/domain"

func HandleEvent(event domain.Event) (domain.Stock, error) {
	return domain.Apply(event)
}
