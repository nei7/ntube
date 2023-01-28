package db

import "errors"

var (
	DuplicateKeyValueError = errors.New("ERROR: duplicate key value violates unique constraint")
)
