package dbutil

import "database/sql"

type TxOption func(options *sql.TxOptions) error

func WithTxIsolationLevel(isolationLevel sql.IsolationLevel) func(options *sql.TxOptions) error {
	return func(options *sql.TxOptions) error {
		options.Isolation = isolationLevel
		return nil
	}
}

func WithTxReadOnly(readOnly bool) func(options *sql.TxOptions) error {
	return func(options *sql.TxOptions) error {
		options.ReadOnly = readOnly
		return nil
	}
}
