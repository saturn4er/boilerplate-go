package authservice

import (
	context "context"

	segmentationsvc "github.com/saturn4er/boilerplate-go/example/segmentation/segmentationservice"
	dbutil "github.com/saturn4er/boilerplate-go/lib/dbutil"
	idempotency "github.com/saturn4er/boilerplate-go/lib/idempotency"
	txoutbox "github.com/saturn4er/boilerplate-go/lib/txoutbox"
	// user code 'imports'
	// end user code 'imports'
)

type Storage interface {
	Users() UsersStorage
	SetUserTagCommands() SetUserTagCommandsOutbox
	IdempotencyKeys() idempotency.Storage
	ExecuteInTransaction(ctx context.Context, cb func(ctx context.Context, tx Storage) error) error
	WithAdvisoryLock(ctx context.Context, scope string, lockID int64) error
	// user code 'Storage custom methods'
	// end user code 'Storage custom methods'
}
type UsersStorage interface {
	dbutil.EntityStorage[User, UserFilter]
	// user code 'User metods'
	// end user code 'User metods'
}

// user code 'User definitions'
// end user code 'User definitions'

type SetUserTagCommandsOutbox txoutbox.Outbox[segmentationsvc.SetUserTagCommand]
