// Package rwvfs augments vfs to support write operations.
package rwvfs

import (
	"code.google.com/p/go.tools/godoc/vfs"
	"io"
)

type FileSystem interface {
	vfs.FileSystem
	Opener
}

type Opener interface {
	OpenFile(path string, flag int) (ReadWriteSeekCloser, error)
}

type ReadWriteSeekCloser interface {
	io.Reader
	io.Writer
	io.Seeker
	io.Closer
}
