package model

import (
	"errors"

	"github.com/quincy0/go-kit/core/stores/mon"
)

var (
	ErrNotFound        = mon.ErrNotFound
	ErrInvalidObjectId = errors.New("invalid objectId")
)
