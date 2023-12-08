package idempotency

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

type ProcessedIdempotencyKey struct {
	IdempotencyKey string `gorm:"primaryKey"`
	Handler        string `gorm:"primaryKey"`
	CreatedAt      time.Time
}

func (a ProcessedIdempotencyKey) TableName() string {
	return "idempotency.processed_idempotency_keys"
}

type Storage interface {
	StoreProcessed(ctx context.Context, idempotencyKey, handler string) error
}

type GormStorage struct {
	DB *gorm.DB
}

func (s GormStorage) StoreProcessed(ctx context.Context, idempotencyKey, handler string) error {
	dbModel := &ProcessedIdempotencyKey{
		IdempotencyKey: idempotencyKey,
		Handler:        handler,
		CreatedAt:      time.Now(),
	}

	err := s.DB.WithContext(ctx).Create(dbModel).Error
	if err != nil {
		var pgErr *pgconn.PgError

		if errors.As(err, &pgErr) {
			if pgErr.Code == "23505" {
				return ErrAlreadyProcessed
			}
		}

		return errors.WithStack(err)
	}

	return nil
}
