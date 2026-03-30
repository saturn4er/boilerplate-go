// Package idempotency provides duplicate message detection using PostgreSQL.
//
// [Storage] defines the interface for recording processed idempotency keys:
//
//	type Storage interface {
//	    StoreProcessed(ctx context.Context, idempotencyKey, handler string) error
//	}
//
// [GormStorage] implements this using a unique constraint on (idempotency_key, handler)
// in the idempotency.processed_idempotency_keys table. When a duplicate key is
// encountered, [ErrAlreadyProcessed] is returned.
//
// This is typically used with the transactional outbox pattern to ensure
// event handlers process each message exactly once.
package idempotency
