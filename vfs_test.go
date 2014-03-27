package rwvfs

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"testing"
)

func TestSub(t *testing.T) {
	m := Map(map[string]string{})
	sub := Sub(m, "/sub")

	err := sub.Mkdir("/")
	if err != nil {
		t.Fatal(err)
	}
	testIsDir(t, "sub", m, "/sub")

	f, err := sub.Create("f1")
	f.Close()
	if err != nil {
		t.Fatal(err)
	}
	testIsFile(t, "sub", m, "/sub/f1")

	f, err = sub.Create("/f2")
	f.Close()
	if err != nil {
		t.Fatal(err)
	}
	testIsFile(t, "sub", m, "/sub/f2")

	err = sub.Mkdir("/d1")
	if err != nil {
		t.Fatal(err)
	}
	testIsDir(t, "sub", m, "/sub/d1")

	err = sub.Mkdir("/d2")
	if err != nil {
		t.Fatal(err)
	}
	testIsDir(t, "sub", m, "/sub/d2")
}

func TestRWVFS(t *testing.T) {
	tmpdir, err := ioutil.TempDir("", "rwvfs-test-")
	if err != nil {
		t.Fatal("TempDir", err)
	}
	defer os.RemoveAll(tmpdir)

	tests := []struct {
		fs   FileSystem
		path string
	}{
		{OS(tmpdir), "/foo"},
		{Map(map[string]string{}), "/foo"},
		{Sub(Map(map[string]string{}), "/x"), "/foo"},
	}
	for _, test := range tests {
		testWrite(t, test.fs, test.path)
		testMkdir(t, test.fs)
		testMkdirAll(t, test.fs)
	}
}

func testWrite(t *testing.T, fs FileSystem, path string) {
	label := fmt.Sprintf("%T", fs)

	w, err := fs.Create(path)
	if err != nil {
		t.Fatalf("%s: WriterOpen: %s", label, err)
	}

	input := []byte("qux")
	_, err = w.Write(input)
	if err != nil {
		t.Fatalf("%s: Write: %s", label, err)
	}

	var r io.ReadCloser
	r, err = fs.Open(path)
	if err != nil {
		t.Fatalf("%s: Open: %s", label, err)
	}
	var output []byte
	output, err = ioutil.ReadAll(r)
	if err != nil {
		t.Fatalf("%s: ReadAll: %s", label, err)
	}
	if !bytes.Equal(output, input) {
		t.Errorf("%s: got output %q, want %q", label, output, input)
	}

	err = w.Close()
	if err != nil {
		t.Fatalf("%s: Close: %s", label, err)
	}

	r, err = fs.Open(path)
	if err != nil {
		t.Fatalf("%s: Open: %s", label, err)
	}
	output, err = ioutil.ReadAll(r)
	if err != nil {
		t.Fatalf("%s: ReadAll: %s", label, err)
	}
	if !bytes.Equal(output, input) {
		t.Errorf("%s: got output %q, want %q", label, output, input)
	}
}

func testMkdir(t *testing.T, fs FileSystem) {
	label := fmt.Sprintf("%T", fs)

	if _, ok := fs.(*subFS); ok {
		fs.Mkdir("/")
	}

	err := fs.Mkdir("dir0")
	if err != nil {
		t.Fatalf("%s: Mkdir(dir0): %s", label, err)
	}
	testIsDir(t, label, fs, "dir0")
	testIsDir(t, label, fs, "/dir0")

	err = fs.Mkdir("/dir1")
	if err != nil {
		t.Fatalf("%s: Mkdir(/dir1): %s", label, err)
	}
	testIsDir(t, label, fs, "dir1")
	testIsDir(t, label, fs, "/dir1")

	err = fs.Mkdir("/dir1")
	if !os.IsExist(err) {
		t.Errorf("%s: Mkdir(/dir1) again: got no error, want os.IsExist-satisfying error", label)
	}

	err = fs.Mkdir("/parent-doesnt-exist/dir2")
	if !os.IsNotExist(err) {
		t.Errorf("%s: Mkdir(/parent-doesnt-exist/dir2): got error %v, want os.IsNotExist-satisfying error", label, err)
	}
}

func testMkdirAll(t *testing.T, fs FileSystem) {
	label := fmt.Sprintf("%T", fs)

	err := MkdirAll(fs, "/a/b/c")
	if err != nil {
		t.Fatalf("%s: MkdirAll: %s", label, err)
	}
	testIsDir(t, label, fs, "/a")
	testIsDir(t, label, fs, "/a/b")
	testIsDir(t, label, fs, "/a/b/c")

	err = MkdirAll(fs, "/a/b/c")
	if err != nil {
		t.Fatalf("%s: MkdirAll again: %s", label, err)
	}
}

func testIsDir(t *testing.T, label string, fs FileSystem, path string) {
	fi, err := fs.Stat(path)
	if err != nil {
		t.Fatalf("%s: Stat(%q): %s", label, path, err)
	}

	if !fi.IsDir() {
		t.Errorf("%s: got fs.Stat(%q) IsDir() == false, want true", label, path)
	}
}

func testIsFile(t *testing.T, label string, fs FileSystem, path string) {
	fi, err := fs.Stat(path)
	if err != nil {
		t.Fatalf("%s: Stat(%q): %s", label, path, err)
	}

	if !fi.Mode().IsRegular() {
		t.Errorf("%s: got fs.Stat(%q) Mode().IsRegular() == false, want true", label, path)
	}
}
