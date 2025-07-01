package teststorage

import (
	errors1 "errors"

	pgconn "github.com/jackc/pgx/v5/pgconn"
	errors "github.com/pkg/errors"
	gorm "gorm.io/gorm"

	testsvc "github.com/saturn4er/boilerplate-go/test/test/testservice"
	// user code 'imports'
	// end user code 'imports'
)

func wrapSomeModelQueryError(err error) error {
	if err == nil {
		return nil
	}

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return errors.WithStack(errors1.Join(testsvc.ErrSomeModelNotFound, err))
	}

	var pgErr *pgconn.PgError

	if errors.As(err, &pgErr) {
		if pgErr.Code == "23505" {
			return errors.WithStack(errors1.Join(testsvc.ErrSomeModelAlreadyExists, err))
		}
	}

	return err
}
func wrapSomeOtherModelQueryError(err error) error {
	if err == nil {
		return nil
	}

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return errors.WithStack(errors1.Join(testsvc.ErrSomeOtherModelNotFound, err))
	}

	var pgErr *pgconn.PgError

	if errors.As(err, &pgErr) {
		if pgErr.Code == "23505" {
			return errors.WithStack(errors1.Join(testsvc.ErrSomeOtherModelAlreadyExists, err))
		}
	}

	return err
}

func wrapPasswordRecoveryEventQueryError(err error) error {
	if err == nil {
		return nil
	}

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return errors.WithStack(errors1.Join(testsvc.ErrPasswordRecoveryEventNotFound, err))
	}

	var pgErr *pgconn.PgError

	if errors.As(err, &pgErr) {
		if pgErr.Code == "23505" {
			return errors.WithStack(errors1.Join(testsvc.ErrPasswordRecoveryEventAlreadyExists, err))
		}
	}

	return err
}
