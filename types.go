package goarpcsolution

import (
	"io"
)

type ReadWriteSeekCloser interface {
	io.Reader
	io.Writer
	io.Closer
	io.Seeker
}
