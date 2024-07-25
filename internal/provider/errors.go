package provider

import "errors"

var (
	ErrFileNotFound = errors.New("no such file or directory")
	ErrInvalidLine  = errors.New("invalid line")
)
