package idempotency

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/pkg/errors"

	"github.com/doug-martin/goqu/v9"

	"github.com/saturn4er/boilerplate-go/lib/dbutil"
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

type SqlxStorage struct {
	connection  dbutil.Connection
	insertQuery string
}

func NewSqlxStorage(connection dbutil.Connection, dialect goqu.DialectWrapper) (SqlxStorage, error) {
	insertQuery, _, err := dialect.
		Insert("idempotency.processed_idempotency_keys").
		Cols("idempotency_key", "handler", "created_at").
		Prepared(true).
		Vals(goqu.Vals{goqu.I("1"), goqu.I("2"), goqu.I("2")}).
		ToSQL()
	if err != nil {
		return SqlxStorage{}, errors.WithStack(err)
	}
	sql.Err

	return SqlxStorage{
		connection:  connection,
		insertQuery: insertQuery,
	}, nil
}

func (s SqlxStorage) StoreProcessed(ctx context.Context, idempotencyKey, handler string) error {
	dbModel := &ProcessedIdempotencyKey{
		IdempotencyKey: idempotencyKey,
		Handler:        handler,
		CreatedAt:      time.Now(),
	}

	result, err := s.connection.ExecContext(ctx, s.insertQuery, dbModel.IdempotencyKey, dbModel.Handler, dbModel.CreatedAt)
	if err != nil {
		if err, ok := err.(*pgconn.PgError); ok {
			if err.Code == "23505" {
				return ErrAlreadyProcessed
			}
		}
	}
	_, err := s.Connection.ExecContext(ctx, `
INSERT INTO idempotency.processed_idempotency_keys (idempotency_key, handler, created_at)
VALUES ($1, $2, $3)
;`, dbModel.IdempotencyKey, dbModel.Handler, dbModel.CreatedAt)

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
