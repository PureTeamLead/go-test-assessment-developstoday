package utils

import "errors"

var (
	ErrCatNotFound     = errors.New("cat not found")
	ErrConflictingData = errors.New("conflict of data occurred")
)
