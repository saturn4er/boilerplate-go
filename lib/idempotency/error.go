package idempotency

import "github.com/pkg/errors"

var ErrAlreadyProcessed = errors.New("already processed")
