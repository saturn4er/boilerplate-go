package txoutbox

import (
	"context"
	"sync/atomic"
	"time"

	"github.com/go-pnp/go-pnp/logging"
	"github.com/go-pnp/jobber"
	"github.com/pkg/errors"
	"gorm.io/gorm"

	"github.com/saturn4er/boilerplate-go/lib/dbutil"
	"github.com/saturn4er/boilerplate-go/lib/filter"
	"github.com/saturn4er/boilerplate-go/lib/pagination"
)

type MessageSender interface {
	SendMessage(ctx context.Context, message *Message) error
}

var _ jobber.Job = new(MessagesSender)

type MessagesSender struct {
	db            *gorm.DB
	logger        *logging.Logger
	messageSender MessageSender

	// count of messages handled in the last Handle method call
	lastHandlingMessagesCount atomic.Int32
}

func NewMessagesProcessor(
	db *gorm.DB,
	logger *logging.Logger,
	messageProcessor MessageSender,
) *MessagesSender {
	return &MessagesSender{
		db:            db,
		logger:        logger,
		messageSender: messageProcessor,
	}
}

func (m *MessagesSender) Init(ctx context.Context) error {
	return nil
}

func (m *MessagesSender) Handle(ctx context.Context) error {
	storage := GormStorage[any]{
		DB: m.db,
	}
	messages, err := storage.Find(ctx, nil, dbutil.WithOrder(MessageFieldCreatedAt, dbutil.OrderDirAsc), dbutil.WithPagination(&pagination.Pagination{
		Page:    1,
		PerPage: 100,
	}))
	if err != nil {
		return errors.WithStack(err)
	}

	m.lastHandlingMessagesCount.Store(int32(len(messages)))

	for _, message := range messages {
		if err := m.db.Transaction(func(tx *gorm.DB) error {
			storage := GormStorage[any]{DB: tx}
			message, err := storage.First(ctx, &MessageFilter{
				ID: filter.Equals(message.ID),
			}, dbutil.WithForUpdate())
			if err != nil {
				if errors.Is(err, gorm.ErrRecordNotFound) {
					return nil
				}

				return err
			}

			if err := m.messageSender.SendMessage(ctx, message); err != nil {
				return err
			}

			return storage.Delete(ctx, &MessageFilter{
				ID: filter.Equals(message.ID),
			})
		}); err != nil {
			m.logger.WithError(err).WithField("message", message).Error(ctx, "failed to send message")
		} else {
			m.logger.WithField("message", message).Info(ctx, "message sent")
		}
	}

	return nil
}

func (m *MessagesSender) Timer() *time.Timer {
	return time.NewTimer(0)
}

func (m *MessagesSender) ResetTimer(timer *time.Timer) {
	if m.lastHandlingMessagesCount.Load() < 30 {
		timer.Reset(time.Second)
	} else {
		timer.Reset(0)
	}
}
