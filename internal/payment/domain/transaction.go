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
	if err != nil {
		return b, nil, err
	}

	return balance, payment, nil
}
