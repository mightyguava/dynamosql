package errors

import (
	"errors"
)

var ErrPositionalArg = errors.New("positional args not supported, use sql.NamedArg to pass named arguments")
