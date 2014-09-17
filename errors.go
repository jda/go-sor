package sor

import (
	"errors"
)

var ErrNoBlock = errors.New("sor: no such block")
var ErrIncompleteBlock = errors.New("sor: incomplete block")
