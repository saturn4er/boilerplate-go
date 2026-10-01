-- +goose NO TRANSACTION
-- +goose Up
CREATE INDEX CONCURRENTLY IF NOT EXISTS processed_idempotency_keys_created_at_idx
    ON idempotency.processed_idempotency_keys (created_at);

-- +goose Down
DROP INDEX CONCURRENTLY IF EXISTS idempotency.processed_idempotency_keys_created_at_idx;
