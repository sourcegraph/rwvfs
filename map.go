package rwvfs

import (
	"io"
	"strings"

	"code.google.com/p/go.tools/godoc/vfs"
	"code.google.com/p/go.tools/godoc/vfs/mapfs"
)

// Map returns a new FileSystem from the provided map. Map keys should be
// forward slash-separated pathnames and not contain a leading slash.
func Map(m map[string]string) FileSystem {
	return mapFS{m, mapfs.New(m)}
}

type mapFS struct {
	m map[string]string
	vfs.FileSystem
}

func (mfs mapFS) Create(path string) (io.WriteCloser, error) {
	// Mimic behavior of OS filesystem: truncate to empty string upon creation;
	// immediately update string values with writes.
	path = filename(path)
	mfs.m[path] = ""
	return &mapFile{mfs.m, path}, nil
}

func filename(p string) string {
	return strings.TrimPrefix(p, "/")
}

type mapFile struct {
	m    map[string]string
	path string
}

func (f *mapFile) Write(p []byte) (int, error) {
	f.m[f.path] = f.m[f.path] + string(p)
	return len(p), nil
}

func (f mapFile) Close() error { return nil }
