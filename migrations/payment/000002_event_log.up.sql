CREATE TABLE IF NOT EXISTS event_log (
	id         SERIAL,
	payload    JSONB  NOT NULL,
	PRIMARY KEY(id)
);

CREATE TABLE IF NOT EXISTS event_offset(
	offset_acked BIGINT
);

INSERT INTO event_offset(offset_acked) VALUES (0);
