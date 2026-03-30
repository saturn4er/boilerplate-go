package segmentationstorage

import (
	errors1 "errors"

	pgconn "github.com/jackc/pgx/v5/pgconn"
	errors "github.com/pkg/errors"
	gorm "gorm.io/gorm"

	segmentationsvc "github.com/saturn4er/boilerplate-go/example/segmentation/segmentationservice"
	// user code 'imports'
	// end user code 'imports'
)

func wrapUserTagQueryError(err error) error {
	if err == nil {
		return nil
	}

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return errors.WithStack(errors1.Join(segmentationsvc.ErrUserTagNotFound, err))
	}

	var pgErr *pgconn.PgError

	if errors.As(err, &pgErr) {
		if pgErr.Code == "23505" {
			return errors.WithStack(errors1.Join(segmentationsvc.ErrUserTagAlreadyExists, err))
		}
	}

	return err
}
func wrapSetUserTagCommandQueryError(err error) error {
	if err == nil {
		return nil
	}

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return errors.WithStack(errors1.Join(segmentationsvc.ErrSetUserTagCommandNotFound, err))
	}

	var pgErr *pgconn.PgError

	if errors.As(err, &pgErr) {
		if pgErr.Code == "23505" {
			return errors.WithStack(errors1.Join(segmentationsvc.ErrSetUserTagCommandAlreadyExists, err))
		}
	}

	return err
}
