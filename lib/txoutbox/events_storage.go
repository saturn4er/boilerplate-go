package txoutbox

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/go-pnp/go-pnp/pkg/optionutil"
	"github.com/pkg/errors"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"github.com/saturn4er/boilerplate-go/lib/dbutil"
	"github.com/saturn4er/boilerplate-go/lib/filter"
)

type dbMessage struct {
	ID             int    `gorm:"primaryKey;autoIncrement:true"`
	Topic          string `gorm:"type:varchar(100);not null"`
	OrderingKey    string `gorm:"type:varchar(100);not null"`
	IdempotencyKey string `gorm:"type:varchar(100);not null"`
	Data           []byte `gorm:"type:bytea;not null"`
	Metadata       string `gorm:"type:jsonb;not null"`
	CreatedAt      time.Time
}

func (d *dbMessage) TableName() string {
	return "tx_outbox.messages"
}

type Message struct {
	ID             int               `gorm:"primaryKey;autoIncrement:true"`
	Topic          string            `gorm:"type:varchar(100);not null"`
	OrderingKey    string            `gorm:"type:varchar(100);not null"`
	IdempotencyKey string            `gorm:"type:varchar(100);not null"`
	Data           []byte            `gorm:"type:bytea;not null"`
	Metadata       map[string]string `gorm:"type:jsonb;not null"`
	CreatedAt      time.Time
}

func convertMessageToDB(src *Message) (*dbMessage, error) {
	metadata, err := json.Marshal(src.Metadata)
	if err != nil {
		return nil, fmt.Errorf("encode metadata: %w", err)
	}

	return &dbMessage{
		ID:             src.ID,
		Topic:          src.Topic,
		OrderingKey:    src.OrderingKey,
		IdempotencyKey: src.IdempotencyKey,
		Data:           src.Data,
		Metadata:       string(metadata),
		CreatedAt:      src.CreatedAt,
	}, nil
}

func convertMessageFromDB(src *dbMessage) (*Message, error) {
	var metadata map[string]string
	err := json.Unmarshal([]byte(src.Metadata), &metadata)
	if err != nil {
		return nil, fmt.Errorf("decode metadata: %w", err)
	}

	return &Message{
		ID:             src.ID,
		Topic:          src.Topic,
		OrderingKey:    src.OrderingKey,
		IdempotencyKey: src.IdempotencyKey,
		Data:           src.Data,
		Metadata:       metadata,
		CreatedAt:      src.CreatedAt,
	}, nil
}

type MessageField byte

const (
	MessageFieldID MessageField = iota + 1
	MessageFieldTopic
	MessageFieldOrderingKey
	MessageFieldIdempotencyKey
	MessageFieldMetadata
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
	message, err := s.BuildMessage(model)
	if err != nil {
		return err
	}

	dbMessageToCreate, err := convertMessageToDB(message)
	if err != nil {
		return err
	}

	err = s.DB.WithContext(ctx).Create(dbMessageToCreate).Error
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
	resultDBMessage := new(dbMessage)
	db := s.DB.WithContext(ctx).Model(&resultDBMessage)

	clauses, err := optionutil.ApplyOptions(&dbutil.SelectOptions{}, options...).BuildExpressions(map[any]clause.Column{
		MessageFieldID:             {Name: "id"},
		MessageFieldTopic:          {Name: "topic"},
		MessageFieldOrderingKey:    {Name: "ordering_key"},
		MessageFieldIdempotencyKey: {Name: "idempotency_key"},
		MessageFieldMetadata:       {Name: "metadata"},
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

	if err := db.First(&resultDBMessage).Error; err != nil {
		return nil, errors.WithStack(err)
	}

	message, err := convertMessageFromDB(resultDBMessage)
	if err != nil {
		return nil, err
	}

	return message, nil
}

func (s GormStorage[ExtType]) Find(
	ctx context.Context,
	filter *MessageFilter,
	options ...optionutil.Option[dbutil.SelectOptions],
) ([]*Message, error) {
	var dbMessages []*dbMessage
	db := s.DB.WithContext(ctx).Model(&dbMessages)

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

	if err := db.Find(&dbMessages).Error; err != nil {
		return nil, errors.WithStack(err)
	}

	messages := make([]*Message, 0, len(dbMessages))

	for _, d := range dbMessages {
		message, err := convertMessageFromDB(d)
		if err != nil {
			return nil, err
		}

		messages = append(messages, message)
	}

	return messages, nil
}

func (s GormStorage[ExtType]) Delete(ctx context.Context, filter *MessageFilter) error {
	var dbTypes []*dbMessage
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
