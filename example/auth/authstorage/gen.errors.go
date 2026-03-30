package authstorage

import (
	errors1 "errors"

	pgconn "github.com/jackc/pgx/v5/pgconn"
	errors "github.com/pkg/errors"
	gorm "gorm.io/gorm"

	authsvc "github.com/saturn4er/boilerplate-go/example/auth/authservice"
	// user code 'imports'
	// end user code 'imports'
)

func wrapUserQueryError(err error) error {
	if err == nil {
		return nil
	}

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return errors.WithStack(errors1.Join(authsvc.ErrUserNotFound, err))
	}

	var pgErr *pgconn.PgError

	if errors.As(err, &pgErr) {
		if pgErr.Code == "23505" {
			return errors.WithStack(errors1.Join(authsvc.ErrUserAlreadyExists, err))
		}
	}

	return err
}
