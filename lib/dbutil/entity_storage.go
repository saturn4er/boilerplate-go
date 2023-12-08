package dbutil

import (
	"context"

	"github.com/go-pnp/go-pnp/logging"
	"github.com/go-pnp/go-pnp/pkg/optionutil"
	"github.com/pkg/errors"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type EntityStorage[Entity, Filter any] interface {
	Create(ctx context.Context, model *Entity) (*Entity, error)
	Count(ctx context.Context, filter *Filter) (int, error)
	Update(ctx context.Context, model *Entity) (*Entity, error)
	Save(ctx context.Context, model *Entity) (*Entity, error)
	First(ctx context.Context, filter *Filter, options ...optionutil.Option[SelectOptions]) (*Entity, error)
	FirstOrCreate(ctx context.Context, filter *Filter, model *Entity, options ...optionutil.Option[SelectOptions]) (*Entity, error)
	Find(ctx context.Context, filter *Filter, options ...optionutil.Option[SelectOptions]) ([]*Entity, error)
	Delete(ctx context.Context, filter *Filter) error
}

type GormEntityStorage[ExtType any, IntType any, FilterType any] struct {
	Logger                *logging.Logger
	DB                    *gorm.DB
	DBErrorsWrapper       func(error) error
	ConvertToInternal     func(*ExtType) (*IntType, error)
	ConvertToExternal     func(*IntType) (*ExtType, error)
	BuildFilterExpression func(*FilterType) (clause.Expression, error)
	FieldMapping          map[any]clause.Column
}

func (s GormEntityStorage[ExtType, IntType, FilterType]) Create(ctx context.Context, model *ExtType) (*ExtType, error) {
	dbModel, err := s.ConvertToInternal(model)
	if err != nil {
		return nil, err
	}

	err = s.DB.WithContext(ctx).Create(dbModel).Error
	if err != nil {
		return nil, s.DBErrorsWrapper(errors.WithStack(err))
	}

	return s.ConvertToExternal(dbModel)
}

func (s GormEntityStorage[ExtType, IntType, FilterType]) Count(ctx context.Context, filter *FilterType) (int, error) {
	expr, err := s.BuildFilterExpression(filter)
	if err != nil {
		return 0, err
	}

	var (
		count     int64
		modelType IntType
	)

	if err := s.DB.WithContext(ctx).Model(modelType).Where(expr).Count(&count).Error; err != nil {
		return 0, s.DBErrorsWrapper(errors.WithStack(err))
	}

	return int(count), nil
}

func (s GormEntityStorage[ExtType, IntType, FilterType]) Update(ctx context.Context, model *ExtType) (*ExtType, error) {
	dbModel, err := s.ConvertToInternal(model)
	if err != nil {
		return nil, err
	}

	err = s.DB.WithContext(ctx).Save(dbModel).Error
	if err != nil {
		return nil, s.DBErrorsWrapper(errors.WithStack(err))
	}

	return s.ConvertToExternal(dbModel)
}

func (s GormEntityStorage[ExtType, IntType, FilterType]) Save(ctx context.Context, model *ExtType) (*ExtType, error) {
	dbModel, err := s.ConvertToInternal(model)
	if err != nil {
		return nil, err
	}

	err = s.DB.WithContext(ctx).Save(dbModel).Error
	if err != nil {
		return nil, s.DBErrorsWrapper(errors.WithStack(err))
	}

	return s.ConvertToExternal(dbModel)
}

func (s GormEntityStorage[ExtType, IntType, FilterType]) First(
	ctx context.Context,
	filter *FilterType,
	options ...optionutil.Option[SelectOptions],
) (*ExtType, error) {
	expr, err := s.BuildFilterExpression(filter)
	if err != nil {
		return nil, err
	}

	result := new(IntType)
	db := s.DB.WithContext(ctx).Model(result)

	clauses, err := optionutil.ApplyOptions(&SelectOptions{}, options...).BuildExpressions(s.FieldMapping)
	if err != nil {
		return nil, err
	}

	db = db.Clauses(clauses...)

	if err := db.First(result, expr).Error; err != nil {
		return nil, s.DBErrorsWrapper(errors.WithStack(err))
	}

	return s.ConvertToExternal(result)
}

// FirstOrCreate returns first record that matches given conditions or creates new one with given values. It's safe to call it in transaction as it does not
// produce error if record does not exist.
func (s GormEntityStorage[ExtType, IntType, FilterType]) FirstOrCreate(
	ctx context.Context,
	filter *FilterType,
	model *ExtType,
	options ...optionutil.Option[SelectOptions],
) (*ExtType, error) {
	count, err := s.Count(ctx, filter)
	if err != nil {
		return nil, err
	}

	if count == 0 {
		if _, err := s.Create(ctx, model); err != nil {
			return nil, err
		}
	}

	return s.First(ctx, filter, options...)
}

func (s GormEntityStorage[ExtType, IntType, FilterType]) Find(
	ctx context.Context,
	filter *FilterType,
	options ...optionutil.Option[SelectOptions],
) ([]*ExtType, error) {
	filterExpr, err := s.BuildFilterExpression(filter)
	if err != nil {
		return nil, err
	}

	var dbTypes []IntType
	db := s.DB.WithContext(ctx).Model(&dbTypes)

	clauses, err := optionutil.ApplyOptions(&SelectOptions{}, options...).BuildExpressions(s.FieldMapping)
	if err != nil {
		return nil, err
	}

	db = db.Clauses(clauses...)

	if err := db.Find(&dbTypes, filterExpr).Error; err != nil {
		return nil, s.DBErrorsWrapper(errors.WithStack(err))
	}

	result := make([]*ExtType, 0, len(dbTypes))

	for _, dbType := range dbTypes {
		authVal, err := s.ConvertToExternal(&dbType) //nolint:gosec
		if err != nil {
			return nil, err
		}

		result = append(result, authVal)
	}

	return result, nil
}

func (s GormEntityStorage[ExtType, IntType, FilterType]) Delete(ctx context.Context, filter *FilterType) error {
	filterExpr, err := s.BuildFilterExpression(filter)
	if err != nil {
		return err
	}

	var dbTypes []*IntType
	db := s.DB.WithContext(ctx).Model(&dbTypes)

	if err := db.Delete(&dbTypes, filterExpr).Error; err != nil {
		return s.DBErrorsWrapper(errors.WithStack(err))
	}

	return nil
}
