// Package rwvfs augments vfs to support write operations.
package rwvfs

import (
	"code.google.com/p/go.tools/godoc/vfs"
	"io"
)

type FileSystem interface {
	vfs.FileSystem
	WriterOpener
}

type WriterOpener interface {
	WriterOpen(path string) (io.WriteCloser, error)
}
