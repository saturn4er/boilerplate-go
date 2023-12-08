-- +goose Up
-- +goose StatementBegin
CREATE SCHEMA IF NOT EXISTS tx_outbox;
CREATE TABLE tx_outbox.messages
(
    id              serial                   NOT NULL,
    topic           varchar(100)             NOT NULL,
    ordering_key    varchar(100)             NOT NULL,
    idempotency_key varchar(100)             NOT NULL,
    data            bytea                    NOT NULL,
    created_at      timestamp with time zone NOT NULL,
    PRIMARY KEY (id)
);
-- index by created at ASC
CREATE INDEX messages_created_at_idx ON tx_outbox.messages (created_at ASC);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE tx_outbox.messages;
DROP SCHEMA tx_outbox;
-- +goose StatementEnd
