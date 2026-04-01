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
	SomeModels() SomeModelsStorage
	SomeOtherModels() SomeOtherModelsStorage
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
type SomeModelsStorage interface {
	dbutil.EntityStorage[SomeModel, SomeModelFilter]
	// user code 'SomeModel metods'
	// end user code 'SomeModel metods'
}

// user code 'SomeModel definitions'
// end user code 'SomeModel definitions'
type SomeOtherModelsStorage interface {
	dbutil.EntityStorage[SomeOtherModel, SomeOtherModelFilter]
	// user code 'SomeOtherModel metods'
	// end user code 'SomeOtherModel metods'
}

// user code 'SomeOtherModel definitions'
// end user code 'SomeOtherModel definitions'
