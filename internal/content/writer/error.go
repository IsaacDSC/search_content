package writer

import "errors"

var ErrAlreadyRegistered = errors.New("endpoint already registered")
var ErrInvalidDataType = errors.New("invalid data type")
