package testservice

import (
	context "context"

	dbutil "github.com/saturn4er/boilerplate-go/lib/dbutil"
	idempotency "github.com/saturn4er/boilerplate-go/lib/idempotency"
	txoutbox "github.com/saturn4er/boilerplate-go/lib/txoutbox"
	// user code 'imports'
	// end user code 'imports'
)

type Storage interface {
	SomeModels() SomeModelsStorage
	SomeOtherModels() SomeOtherModelsStorage
	PasswordRecoveryEvents() PasswordRecoveryEventsOutbox
	IdempotencyKeys() idempotency.Storage
	ExecuteInTransaction(ctx context.Context, cb func(ctx context.Context, tx Storage) error) error
}
type SomeModelsStorage dbutil.EntityStorage[SomeModel, SomeModelFilter]
type SomeOtherModelsStorage dbutil.EntityStorage[SomeOtherModel, SomeOtherModelFilter]
type PasswordRecoveryEventsOutbox txoutbox.Outbox[PasswordRecoveryEvent]
