CREATE TYPE payment_status AS ENUM ('new', 'completed', 'canceled');

CREATE TABLE IF NOT EXISTS balances (
	customer_id      UUID UNIQUE NOT NULL,
	available_amount DECIMAL     NOT NULL DEFAULT 0,
	reserved_amount  DECIMAL     NOT NULL DEFAULT 0,
	PRIMARY KEY(customer_id)
);

CREATE TABLE IF NOT EXISTS payments (
	payment_id  UUID UNIQUE    NOT NULL,
	customer_id UUID           NOT NULL,
	order_id    UUID UNIQUE    NOT NULL,
	amount      DECIMAL        NOT NULL,
	status      payment_status NOT NULL DEFAULT 'new',
	PRIMARY KEY(payment_id),
	CONSTRAINT fk_customer FOREIGN KEY (customer_id)
		REFERENCES balances(customer_id)
);
