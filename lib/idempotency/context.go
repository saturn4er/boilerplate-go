package idempotency

import "context"

type idempotencyKeyCtxKey struct{}

func ContextWithIdempotencyKey(ctx context.Context, idempotencyKey string) context.Context {
	return context.WithValue(ctx, idempotencyKeyCtxKey{}, idempotencyKey)
}

func IdempotencyKeyFromContext(ctx context.Context) (string, bool) {
	idempotencyKey, ok := ctx.Value(idempotencyKeyCtxKey{}).(string)

	return idempotencyKey, ok
}
