package authstorage

import (
	context "context"
	strconv "strconv"

	xxhash "github.com/cespare/xxhash"
	logging "github.com/go-pnp/go-pnp/logging"
	gorm "gorm.io/gorm"
	clause "gorm.io/gorm/clause"

	authsvc "github.com/saturn4er/boilerplate-go/example/auth/authservice"
	segmentationevent "github.com/saturn4er/boilerplate-go/example/segmentation/segmentationevent"
	segmentationsvc "github.com/saturn4er/boilerplate-go/example/segmentation/segmentationservice"
	dbutil "github.com/saturn4er/boilerplate-go/lib/dbutil"
	idempotency "github.com/saturn4er/boilerplate-go/lib/idempotency"
	txoutbox "github.com/saturn4er/boilerplate-go/lib/txoutbox"
	// user code 'imports'
	// end user code 'imports'
)

type Storages struct {
	db         *gorm.DB
	logger     *logging.Logger
	processors []txoutbox.MessageProcessor
}

var _ authsvc.Storage = &Storages{}

func (s Storages) Users() authsvc.UsersStorage {
	return NewUsersStorage(s.db, s.logger)
}

func (s Storages) SetUserTagCommands() authsvc.SetUserTagCommandsOutbox {
	return NewSetUserTagCommandsOutbox(s.db, s.processors)
}

func (s Storages) IdempotencyKeys() idempotency.Storage {
	return idempotency.GormStorage{
		DB: s.db,
	}

}

func (s *Storages) WithAdvisoryLock(ctx context.Context, scope string, lockID int64) error {
	hasher := xxhash.New()
	hasher.Write([]byte(scope))
	hasher.Write([]byte{':'})
	hasher.Write(strconv.AppendInt(nil, lockID, 10))

	result := s.db.WithContext(ctx).Exec("SELECT pg_advisory_xact_lock(?)", int64(hasher.Sum64()))
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (s Storages) ExecuteInTransaction(ctx context.Context, cb func(ctx context.Context, tx authsvc.Storage) error) error {
	return s.db.Transaction(func(tx *gorm.DB) error {
		return cb(ctx, &Storages{db: tx, logger: s.logger, processors: s.processors})
	})
}

func NewStorages(db *gorm.DB, logger *logging.Logger, processors []txoutbox.MessageProcessor) *Storages {
	return &Storages{db: db, logger: logger, processors: processors}
}

type UsersStorage struct {
	dbutil.GormEntityStorage[authsvc.User, dbUser, authsvc.UserFilter]
}

// user code 'User custom methods'
// end user code 'User custom methods'
func NewUsersStorage(db *gorm.DB, logger *logging.Logger) authsvc.UsersStorage {
	return &UsersStorage{
		GormEntityStorage: dbutil.GormEntityStorage[authsvc.User, dbUser, authsvc.UserFilter]{
			Logger:            logger,
			DB:                db,
			DBErrorsWrapper:   wrapUserQueryError,
			ConvertToInternal: convertUserToDB,
			ConvertToExternal: convertUserFromDB,
			BuildFilterExpression: func(filter *authsvc.UserFilter) (clause.Expression, error) {
				return buildUserFilterExpr(filter)
			},
			FieldMapping: map[any]clause.Column{
				authsvc.UserFieldID:    {Name: "id"},
				authsvc.UserFieldEmail: {Name: "email"},
				authsvc.UserFieldName:  {Name: "name"},
				authsvc.UserFieldRole:  {Name: "role"},
			},
			LockScope: "auth.Users",
		},
	}
}

func NewSetUserTagCommandsOutbox(db *gorm.DB, processors []txoutbox.MessageProcessor) authsvc.SetUserTagCommandsOutbox {
	return txoutbox.GormStorage[segmentationsvc.SetUserTagCommand]{
		DB: db,

		BuildMessage:      segmentationevent.BuildSetUserTagCommandMessage,
		MessageProcessors: processors,
	}
}
