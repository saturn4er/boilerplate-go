package txoutbox

import (
	"context"
	"time"

	"github.com/go-pnp/go-pnp/pkg/optionutil"
	"github.com/pkg/errors"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"github.com/saturn4er/boilerplate-go/lib/dbutil"
	"github.com/saturn4er/boilerplate-go/lib/filter"
)

type Message struct {
	ID             int    `gorm:"primaryKey;autoIncrement:true"`
	Topic          string `gorm:"type:varchar(100);not null"`
	OrderingKey    string `gorm:"type:varchar(100);not null"`
	IdempotencyKey string `gorm:"type:varchar(100);not null"`
	Data           []byte `gorm:"type:bytea;not null"`
	CreatedAt      time.Time
}

func (a Message) TableName() string {
	return "tx_outbox.messages"
}

type MessageField byte

const (
	MessageFieldID MessageField = iota + 1
	MessageFieldTopic
	MessageFieldOrderingKey
	MessageFieldIdempotencyKey
	MessageFieldCreatedAt
)

type MessageFilter struct {
	ID    filter.Filter[int]
	Topic filter.Filter[string]
}

func (e *MessageFilter) buildExpression() (clause.Expression, error) {
	if e == nil {
		return nil, nil
	}
	return dbutil.BuildFilterExpression(
		dbutil.ColumnFilter[int]{Column: "id", Filter: e.ID},
		dbutil.ColumnFilter[string]{Column: "topic", Filter: e.Topic},
	)
}

type Outbox[Entity any] interface {
	Send(ctx context.Context, model *Entity) error
}

type GormStorage[ExtType any] struct {
	DB           *gorm.DB
	BuildMessage func(*ExtType) (*Message, error)
}

func (s GormStorage[ExtType]) Send(ctx context.Context, model *ExtType) error {
	dbModel, err := s.BuildMessage(model)
	if err != nil {
		return err
	}

	err = s.DB.WithContext(ctx).Create(dbModel).Error
	if err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func (s GormStorage[ExtType]) First(
	ctx context.Context,
	filter *MessageFilter,
	options ...optionutil.Option[dbutil.SelectOptions],
) (*Message, error) {
	result := new(Message)
	db := s.DB.WithContext(ctx).Model(&result)

	clauses, err := optionutil.ApplyOptions(&dbutil.SelectOptions{}, options...).BuildExpressions(map[any]clause.Column{
		MessageFieldID:             {Name: "id"},
		MessageFieldTopic:          {Name: "topic"},
		MessageFieldOrderingKey:    {Name: "ordering_key"},
		MessageFieldIdempotencyKey: {Name: "idempotency_key"},
		MessageFieldCreatedAt:      {Name: "created_at"},
	})
	if err != nil {
		return nil, err
	}

	db = db.Clauses(clauses...)

	filterExpr, err := filter.buildExpression()
	if err != nil {
		return nil, err
	}

	if filterExpr != nil {
		db = db.Clauses(filterExpr)
	}

	if err := db.First(&result).Error; err != nil {
		return nil, errors.WithStack(err)
	}

	return result, nil
}

func (s GormStorage[ExtType]) Find(
	ctx context.Context,
	filter *MessageFilter,
	options ...optionutil.Option[dbutil.SelectOptions],
) ([]*Message, error) {
	var result []*Message
	db := s.DB.WithContext(ctx).Model(&result)

	clauses, err := optionutil.ApplyOptions(&dbutil.SelectOptions{}, options...).BuildExpressions(map[any]clause.Column{
		MessageFieldID:             {Name: "id"},
		MessageFieldTopic:          {Name: "topic"},
		MessageFieldOrderingKey:    {Name: "ordering_key"},
		MessageFieldIdempotencyKey: {Name: "idempotency_key"},
		MessageFieldCreatedAt:      {Name: "created_at"},
	})
	if err != nil {
		return nil, err
	}

	db = db.Clauses(clauses...)

	filterExpr, err := filter.buildExpression()
	if err != nil {
		return nil, err
	}

	if filterExpr != nil {
		db = db.Clauses(filterExpr)
	}

	if err := db.Find(&result).Error; err != nil {
		return nil, errors.WithStack(err)
	}

	return result, nil
}

func (s GormStorage[ExtType]) Delete(ctx context.Context, filter *MessageFilter) error {
	var dbTypes []*Message
	db := s.DB.WithContext(ctx).Model(&dbTypes)

	filterExpr, err := filter.buildExpression()
	if err != nil {
		return err
	}

	if filterExpr != nil {
		db = db.Clauses(filterExpr)
	}

	if err := db.Delete(&dbTypes, filterExpr).Error; err != nil {
		return errors.WithStack(err)
	}

	return nil
}
