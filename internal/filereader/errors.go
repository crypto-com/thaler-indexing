package filereader

import (
	"errors"
)

var (
	ErrFileNotFound = errors.New("File not found")
	ErrReadFile     = errors.New("File read error")
)
