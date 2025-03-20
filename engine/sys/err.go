package sys

import "errors"

var ErrNoFound error = errors.New("no found")
var ErrInternalServer error = errors.New("internal server error")
var ErrCategory error = errors.New("category error")
