package rwvfs

import (
	"fmt"
	"os"
	pathpkg "path"
	"path/filepath"

	"code.google.com/p/go.tools/godoc/vfs"
)

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

// OpenFile opens the file at path for writing, creating the file if it doesn't
// exist and truncating it otherwise.
func (fs osFS) OpenFile(path string, flag int) (ReadWriteSeekCloser, error) {
	f, err := os.OpenFile(fs.resolve(path), flag, 0644)
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
