package provider

import "errors"

var ErrFileNotFound = errors.New("no such file or directory")
var ErrInvalidLine = errors.New("invalid line")
