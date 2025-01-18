-- +goose Up
-- +goose StatementBegin
ALTER TABLE tx_outbox.messages
    ADD COLUMN metadata jsonb NOT NULL DEFAULT '{}';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE tx_outbox.messages
    DROP COLUMN metadata;
-- +goose StatementEnd
