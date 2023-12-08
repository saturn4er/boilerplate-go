-- +goose Up
-- +goose StatementBegin
CREATE SCHEMA IF NOT EXISTS idempotency;
CREATE TABLE idempotency.processed_idempotency_keys
(
    idempotency_key VARCHAR(255) NOT NULL,
    handler         VARCHAR(255) NOT NULL,
    created_at      TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (idempotency_key, handler)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE idempotency.processed_idempotency_keys;
DROP SCHEMA idempotency;
-- +goose StatementEnd
