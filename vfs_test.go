package rwvfs

import (
	"bytes"
	"io"
	"io/ioutil"
	"os"
	"testing"
)

func TestWriterOpen(t *testing.T) {
	tmpdir, err := ioutil.TempDir("", "httpstream-test")
	if err != nil {
		t.Fatal("TempDir", err)
	}
	defer os.RemoveAll(tmpdir)

	tests := []struct {
		fs   FileSystem
		path string
	}{
		{OS(tmpdir), "/foo"},
	}
	for _, test := range tests {
		testWriterOpen(t, test.fs, test.path)
	}
}

func testWriterOpen(t *testing.T, fs FileSystem, path string) {
	w, err := fs.WriterOpen(path)
	if err != nil {
		t.Fatal("WriterOpen", err)
	}

	input := []byte("qux")
	_, err = w.Write(input)
	if err != nil {
		t.Fatal("Write", err)
	}

	err = w.Close()
	if err != nil {
		t.Fatal("Close", err)
	}

	var r io.ReadCloser
	r, err = fs.Open(path)
	if err != nil {
		t.Fatal("Open", err)
	}

	var output []byte
	output, err = ioutil.ReadAll(r)
	if err != nil {
		t.Fatal("ReadAll", err)
	}

	if !bytes.Equal(output, input) {
		t.Errorf("got output %q, want %q", output, input)
	}
}
