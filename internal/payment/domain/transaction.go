package domain

type Tx struct {
	Payment Payment
	Event   Event
}

func (b Balance) Transaction(tx Tx) (Balance, Payment, error) {
	payment, err := Apply(tx.Payment, tx.Event)
	if err != nil {
		return b, nil, err
	}

	balance, err := b.Apply(payment)
	switch err {
	case nil:
		return balance, payment, nil
	case ErrInsufficientFunds:
		return b, FailPayment(payment), nil
	default:
		return b, nil, err
	}
}
