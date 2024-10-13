package teststorage

import (
	context "context"

	logging "github.com/go-pnp/go-pnp/logging"
	pgconn "github.com/jackc/pgx/v5/pgconn"
	errors "github.com/pkg/errors"
	gorm "gorm.io/gorm"
	clause "gorm.io/gorm/clause"

	dbutil "github.com/saturn4er/boilerplate-go/lib/dbutil"
	idempotency "github.com/saturn4er/boilerplate-go/lib/idempotency"
	txoutbox "github.com/saturn4er/boilerplate-go/lib/txoutbox"
	testsvc "github.com/saturn4er/boilerplate-go/test/test/testservice"
)

type Storages struct {
	db     *gorm.DB
	logger *logging.Logger
}

var _ testsvc.Storage = &Storages{}

func (s Storages) SomeModels() testsvc.SomeModelsStorage {
	return NewSomeModelsStorage(s.db, s.logger)
}
func (s Storages) SomeOtherModels() testsvc.SomeOtherModelsStorage {
	return NewSomeOtherModelsStorage(s.db, s.logger)
}

func (s Storages) PasswordRecoveryEvents() testsvc.PasswordRecoveryEventsOutbox {
	return NewPasswordRecoveryEventsOutbox(s.db)
}

func (s Storages) IdempotencyKeys() idempotency.Storage {
	return idempotency.GormStorage{
		DB: s.db,
	}
}
func (s Storages) ExecuteInTransaction(ctx context.Context, cb func(ctx context.Context, tx testsvc.Storage) error) error {
	return s.db.Transaction(func(tx *gorm.DB) error {
		return cb(ctx, &Storages{tx, s.logger})
	})
}

func NewStorages(db *gorm.DB, logger *logging.Logger) *Storages {
	return &Storages{db: db, logger: logger}
}

func NewSomeModelsStorage(db *gorm.DB, logger *logging.Logger) testsvc.SomeModelsStorage {
	return dbutil.GormEntityStorage[testsvc.SomeModel, dbSomeModel, testsvc.SomeModelFilter]{
		Logger: logger,
		DB:     db,
		DBErrorsWrapper: func(err error) error {
			if err == nil {
				return nil
			}

			if errors.Is(err, gorm.ErrRecordNotFound) {
				return testsvc.ErrSomeModelNotFound
			}

			var pgErr *pgconn.PgError

			if errors.As(err, &pgErr) {
				if pgErr.Code == "23505" {
					return testsvc.ErrSomeModelAlreadyExists
				}
			}

			// user code 'SomeModel errors mapping'
			// end user code 'SomeModel errors mapping'

			return err
		},
		ConvertToInternal:     convertSomeModelToDB,
		ConvertToExternal:     convertSomeModelFromDB,
		BuildFilterExpression: buildSomeModelFilterExpr,
		FieldMapping: map[any]clause.Column{
			testsvc.SomeModelFieldID:                 {Name: "id"},
			testsvc.SomeModelFieldModelField:         {Name: "model_field"},
			testsvc.SomeModelFieldModelPtrField:      {Name: "model_ptr_field"},
			testsvc.SomeModelFieldOneOfField:         {Name: "one_of_field"},
			testsvc.SomeModelFieldOneOfPtrField:      {Name: "one_of_ptr_field"},
			testsvc.SomeModelFieldEnumField:          {Name: "enum_field"},
			testsvc.SomeModelFieldEnumPtrField:       {Name: "enum_ptr_field"},
			testsvc.SomeModelFieldAnyField:           {Name: "any_field"},
			testsvc.SomeModelFieldAnyPtrField:        {Name: "any_ptr_field"},
			testsvc.SomeModelFieldMapModelField:      {Name: "map_model_field"},
			testsvc.SomeModelFieldMapModelPtrField:   {Name: "map_model_ptr_field"},
			testsvc.SomeModelFieldMapOneOfField:      {Name: "map_one_of_field"},
			testsvc.SomeModelFieldMapOneOfPtrField:   {Name: "map_one_of_ptr_field"},
			testsvc.SomeModelFieldMapEnumField:       {Name: "map_enum_field"},
			testsvc.SomeModelFieldMapEnumPtrField:    {Name: "map_enum_ptr_field"},
			testsvc.SomeModelFieldMapAnyField:        {Name: "map_any_field"},
			testsvc.SomeModelFieldMapAnyPtrField:     {Name: "map_any_ptr_field"},
			testsvc.SomeModelFieldModelSliceField:    {Name: "model_slice_field"},
			testsvc.SomeModelFieldModelPtrSliceField: {Name: "model_ptr_slice_field"},
			testsvc.SomeModelFieldOneOfSliceField:    {Name: "one_of_slice_field"},
			testsvc.SomeModelFieldOneOfPtrSliceField: {Name: "one_of_ptr_slice_field"},
			testsvc.SomeModelFieldSliceEnumField:     {Name: "slice_enum_field"},
			testsvc.SomeModelFieldSliceEnumPtrField:  {Name: "slice_enum_ptr_field"},
			testsvc.SomeModelFieldSliceAnyField:      {Name: "slice_any_field"},
			testsvc.SomeModelFieldSliceAnyPtrField:   {Name: "slice_any_ptr_field"},
		},
	}
}

func NewSomeOtherModelsStorage(db *gorm.DB, logger *logging.Logger) testsvc.SomeOtherModelsStorage {
	return dbutil.GormEntityStorage[testsvc.SomeOtherModel, dbSomeOtherModel, testsvc.SomeOtherModelFilter]{
		Logger: logger,
		DB:     db,
		DBErrorsWrapper: func(err error) error {
			if err == nil {
				return nil
			}

			if errors.Is(err, gorm.ErrRecordNotFound) {
				return testsvc.ErrSomeOtherModelNotFound
			}

			var pgErr *pgconn.PgError

			if errors.As(err, &pgErr) {
				if pgErr.Code == "23505" {
					return testsvc.ErrSomeOtherModelAlreadyExists
				}
			}

			// user code 'SomeOtherModel errors mapping'
			// end user code 'SomeOtherModel errors mapping'

			return err
		},
		ConvertToInternal:     convertSomeOtherModelToDB,
		ConvertToExternal:     convertSomeOtherModelFromDB,
		BuildFilterExpression: buildSomeOtherModelFilterExpr,
		FieldMapping: map[any]clause.Column{
			testsvc.SomeOtherModelFieldID: {Name: "id"},
		},
	}
}

func NewPasswordRecoveryEventsOutbox(db *gorm.DB) testsvc.PasswordRecoveryEventsOutbox {
	return txoutbox.GormStorage[testsvc.PasswordRecoveryEvent]{
		DB:           db,
		BuildMessage: buildPasswordRecoveryEventMessage,
	}
}
