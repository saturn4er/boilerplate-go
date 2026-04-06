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

func (s Storages) SomeModels() segmentationsvc.SomeModelsStorage {
	return NewSomeModelsStorage(s.db, s.logger)
}
func (s Storages) SomeOtherModels() segmentationsvc.SomeOtherModelsStorage {
	return NewSomeOtherModelsStorage(s.db, s.logger)
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
	}
}

type SomeModelsStorage struct {
	dbutil.GormEntityStorage[segmentationsvc.SomeModel, dbSomeModel, segmentationsvc.SomeModelFilter]
}

// user code 'SomeModel custom methods'
// end user code 'SomeModel custom methods'
func NewSomeModelsStorage(db *gorm.DB, logger *logging.Logger) segmentationsvc.SomeModelsStorage {
	return &SomeModelsStorage{
		GormEntityStorage: dbutil.GormEntityStorage[segmentationsvc.SomeModel, dbSomeModel, segmentationsvc.SomeModelFilter]{
			Logger:            logger,
			DB:                db,
			DBErrorsWrapper:   wrapSomeModelQueryError,
			ConvertToInternal: convertSomeModelToDB,
			ConvertToExternal: convertSomeModelFromDB,
			BuildFilterExpression: func(filter *segmentationsvc.SomeModelFilter) (clause.Expression, error) {
				return buildSomeModelFilterExpr(filter)
			},
			FieldMapping: map[any]clause.Column{
				segmentationsvc.SomeModelFieldID:                 {Name: "id"},
				segmentationsvc.SomeModelFieldName:               {Name: "name"},
				segmentationsvc.SomeModelFieldDescription:        {Name: "description"},
				segmentationsvc.SomeModelFieldModelField:         {Name: "model_field"},
				segmentationsvc.SomeModelFieldModelPtrField:      {Name: "model_ptr_field"},
				segmentationsvc.SomeModelFieldOneOfField:         {Name: "one_of_field"},
				segmentationsvc.SomeModelFieldOneOfPtrField:      {Name: "one_of_ptr_field"},
				segmentationsvc.SomeModelFieldEnumField:          {Name: "enum_field"},
				segmentationsvc.SomeModelFieldEnumPtrField:       {Name: "enum_ptr_field"},
				segmentationsvc.SomeModelFieldAnyField:           {Name: "any_field"},
				segmentationsvc.SomeModelFieldAnyPtrField:        {Name: "any_ptr_field"},
				segmentationsvc.SomeModelFieldMapModelField:      {Name: "map_model_field"},
				segmentationsvc.SomeModelFieldMapModelPtrField:   {Name: "map_model_ptr_field"},
				segmentationsvc.SomeModelFieldMapOneOfField:      {Name: "map_one_of_field"},
				segmentationsvc.SomeModelFieldMapOneOfPtrField:   {Name: "map_one_of_ptr_field"},
				segmentationsvc.SomeModelFieldMapEnumField:       {Name: "map_enum_field"},
				segmentationsvc.SomeModelFieldMapEnumPtrField:    {Name: "map_enum_ptr_field"},
				segmentationsvc.SomeModelFieldMapAnyField:        {Name: "map_any_field"},
				segmentationsvc.SomeModelFieldMapAnyPtrField:     {Name: "map_any_ptr_field"},
				segmentationsvc.SomeModelFieldModelSliceField:    {Name: "model_slice_field"},
				segmentationsvc.SomeModelFieldModelPtrSliceField: {Name: "model_ptr_slice_field"},
				segmentationsvc.SomeModelFieldOneOfSliceField:    {Name: "one_of_slice_field"},
				segmentationsvc.SomeModelFieldOneOfPtrSliceField: {Name: "one_of_ptr_slice_field"},
				segmentationsvc.SomeModelFieldSliceEnumField:     {Name: "slice_enum_field"},
				segmentationsvc.SomeModelFieldSliceEnumPtrField:  {Name: "slice_enum_ptr_field"},
				segmentationsvc.SomeModelFieldSliceAnyField:      {Name: "slice_any_field"},
				segmentationsvc.SomeModelFieldSliceAnyPtrField:   {Name: "slice_any_ptr_field"},
			},
			LockScope: "segmentation.SomeModels",
		},
	}
}

type SomeOtherModelsStorage struct {
	dbutil.GormEntityStorage[segmentationsvc.SomeOtherModel, dbSomeOtherModel, segmentationsvc.SomeOtherModelFilter]
}

// user code 'SomeOtherModel custom methods'
// end user code 'SomeOtherModel custom methods'
func NewSomeOtherModelsStorage(db *gorm.DB, logger *logging.Logger) segmentationsvc.SomeOtherModelsStorage {
	return &SomeOtherModelsStorage{
		GormEntityStorage: dbutil.GormEntityStorage[segmentationsvc.SomeOtherModel, dbSomeOtherModel, segmentationsvc.SomeOtherModelFilter]{
			Logger:            logger,
			DB:                db,
			DBErrorsWrapper:   wrapSomeOtherModelQueryError,
			ConvertToInternal: convertSomeOtherModelToDB,
			ConvertToExternal: convertSomeOtherModelFromDB,
			BuildFilterExpression: func(filter *segmentationsvc.SomeOtherModelFilter) (clause.Expression, error) {
				return buildSomeOtherModelFilterExpr(filter)
			},
			FieldMapping: map[any]clause.Column{
				segmentationsvc.SomeOtherModelFieldID: {Name: "id"},
			},
			LockScope: "segmentation.SomeOtherModels",
		},
	}
}
