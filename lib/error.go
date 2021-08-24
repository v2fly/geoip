package lib

import "errors"

var (
	ErrDuplicatedConverter = errors.New("duplicated converter")
	ErrUnknownAction       = errors.New("unknown action")
	ErrInvalidIPType       = errors.New("invalid IP type")
	ErrInvalidIP           = errors.New("invalid IP address")
	ErrInvalidIPLength     = errors.New("invalid IP address length")
	ErrInvalidIPNet        = errors.New("invalid IPNet address")
	ErrInvalidPrefixType   = errors.New("invalid prefix type")
)
