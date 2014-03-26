// Package rwvfs augments vfs to support write operations.
package rwvfs

import (
	"io"

	"code.google.com/p/go.tools/godoc/vfs"
)

type FileSystem interface {
	vfs.FileSystem
	Creator
}

// A Creator implements the Create function for writing to a file.
type Creator interface {
	// Create creates the named file, truncating it if it already exists.
	Create(path string) (ReadWriteSeekCloser, error)
}

type ReadWriteSeekCloser interface {
	io.Reader
	io.Writer
	io.Seeker
	io.Closer
}
