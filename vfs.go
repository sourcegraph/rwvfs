// Package rwvfs augments vfs to support write operations.
package rwvfs

import (
	"bytes"
	"code.google.com/p/go.tools/godoc/vfs"
	"code.google.com/p/go.tools/godoc/vfs/mapfs"
	"fmt"
	"io"
	"os"
	pathpkg "path"
	"path/filepath"
)

type FileSystem interface {
	vfs.FileSystem
	WriterOpener
}

type WriterOpener interface {
	WriterOpen(path string) (io.WriteCloser, error)
}

// OS returns an implementation of FileSystem reading from the tree rooted at
// root.
func OS(root string) FileSystem {
	return osFS{root, vfs.OS(root)}
}

type osFS struct {
	root string
	vfs.FileSystem
}

// resolve is from code.google.com/p/go.tools/godoc/vfs.
func (fs osFS) resolve(path string) string {
	// Clean the path so that it cannot possibly begin with ../.
	// If it did, the result of filepath.Join would be outside the
	// tree rooted at root.  We probably won't ever see a path
	// with .. in it, but be safe anyway.
	path = pathpkg.Clean("/" + path)

	return filepath.Join(string(fs.root), path)
}

// WriterOpen opens the file at path for writing, creating the file if it
// doesn't exist and truncating it otherwise.
func (fs osFS) WriterOpen(path string) (io.WriteCloser, error) {
	f, err := os.OpenFile(fs.resolve(path), os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return nil, err
	}

	fi, err := f.Stat()
	if err != nil {
		return nil, err
	}
	if fi.IsDir() {
		return nil, fmt.Errorf("Open: %s is a directory", path)
	}

	return f, nil
}

// Map returns an implementation of FileSystem based on a map.
func Map(m map[string]string) FileSystem {
	return mapFS{m, mapfs.New(m)}
}

// mapFS is the map based implementation of FileSystem.
type mapFS struct {
	m map[string]string
	vfs.FileSystem
}

// filename is from code.google.com/p/go.tools/godoc/vfs.
func filename(p string) string {
	if len(p) > 0 && p[0] == '/' {
		p = p[1:]
	}
	return p
}

func (fs mapFS) WriterOpen(path string) (io.WriteCloser, error) {
	return &mapWriter{fs: fs, filename: filename(path)}, nil
}

type mapWriter struct {
	fs       mapFS
	filename string
	bytes.Buffer
}

func (w *mapWriter) Write(p []byte) (n int, err error) {
	return w.Buffer.Write(p)
}

func (w *mapWriter) Close() error {
	w.fs.m[w.filename] = w.Buffer.String()
	return nil
}
