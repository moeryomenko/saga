CREATE TYPE order_kind AS ENUM ('empty', 'active', 'pending', 'stocked', 'paid', 'completed', 'canceled');

CREATE TABLE IF NOT EXISTS orders (
	order_id UUID UNIQUE NOT NULL,
	customer_id UUID NOT NULL,
	items JSONB DEFAULT NULL,
	price DECIMAL DEFAULT NULL,
	payment_id UUID DEFAULT NULL,
	kind order_kind NOT NULL DEFAULT 'empty',
	created_at TIMESTAMP(6) WITHOUT TIME ZONE DEFAULT CURRENT_TIMESTAMP,
	updated_at TIMESTAMP(6) WITHOUT TIME ZONE DEFAULT NULL
);
