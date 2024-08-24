package common

import "errors"

var (
	ErrCellIdNotFound     = errors.New("cell id not found")
	ErrDBOperationFailure = errors.New("db operation failure")
)
