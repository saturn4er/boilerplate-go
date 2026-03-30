package segmentationstorage

import (
	context "context"
	strconv "strconv"

	xxhash "github.com/cespare/xxhash"
	logging "github.com/go-pnp/go-pnp/logging"
	gorm "gorm.io/gorm"
	clause "gorm.io/gorm/clause"

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

var _ segmentationsvc.Storage = &Storages{}

func (s Storages) UserTags() segmentationsvc.UserTagsStorage {
	return NewUserTagsStorage(s.db, s.logger)
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

func (s Storages) ExecuteInTransaction(ctx context.Context, cb func(ctx context.Context, tx segmentationsvc.Storage) error) error {
	return s.db.Transaction(func(tx *gorm.DB) error {
		return cb(ctx, &Storages{db: tx, logger: s.logger, processors: s.processors})
	})
}

func NewStorages(db *gorm.DB, logger *logging.Logger, processors []txoutbox.MessageProcessor) *Storages {
	return &Storages{db: db, logger: logger, processors: processors}
}

type UserTagsStorage struct {
	dbutil.GormEntityStorage[segmentationsvc.UserTag, dbUserTag, segmentationsvc.UserTagFilter]
}

// user code 'UserTag custom methods'
// end user code 'UserTag custom methods'
func NewUserTagsStorage(db *gorm.DB, logger *logging.Logger) segmentationsvc.UserTagsStorage {
	return &UserTagsStorage{
		GormEntityStorage: dbutil.GormEntityStorage[segmentationsvc.UserTag, dbUserTag, segmentationsvc.UserTagFilter]{
			Logger:            logger,
			DB:                db,
			DBErrorsWrapper:   wrapUserTagQueryError,
			ConvertToInternal: convertUserTagToDB,
			ConvertToExternal: convertUserTagFromDB,
			BuildFilterExpression: func(filter *segmentationsvc.UserTagFilter) (clause.Expression, error) {
				return buildUserTagFilterExpr(filter)
			},
			FieldMapping: map[any]clause.Column{
				segmentationsvc.UserTagFieldID:     {Name: "id"},
				segmentationsvc.UserTagFieldUserID: {Name: "user_id"},
				segmentationsvc.UserTagFieldKey:    {Name: "key"},
				segmentationsvc.UserTagFieldValue:  {Name: "value"},
			},
			LockScope: "segmentation.UserTags",
		},
		// user code 'UserTag custom metods'
		// end user code 'UserTag custom metods'
	}
}
