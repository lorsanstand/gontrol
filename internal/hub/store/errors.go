package store

import "errors"

var ErrChangeData error = errors.New("failed to change the data")
var ErrNotFound error = errors.New("data not found")
