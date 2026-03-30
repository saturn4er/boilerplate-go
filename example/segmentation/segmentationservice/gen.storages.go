package segmentationservice

import (
	context "context"

	dbutil "github.com/saturn4er/boilerplate-go/lib/dbutil"
	idempotency "github.com/saturn4er/boilerplate-go/lib/idempotency"
	// user code 'imports'
	// end user code 'imports'
)

type Storage interface {
	UserTags() UserTagsStorage
	IdempotencyKeys() idempotency.Storage
	ExecuteInTransaction(ctx context.Context, cb func(ctx context.Context, tx Storage) error) error
	WithAdvisoryLock(ctx context.Context, scope string, lockID int64) error
	// user code 'Storage custom methods'
	// end user code 'Storage custom methods'
}
type UserTagsStorage interface {
	dbutil.EntityStorage[UserTag, UserTagFilter]
	// user code 'UserTag metods'
	// end user code 'UserTag metods'
}

// user code 'UserTag definitions'
// end user code 'UserTag definitions'
