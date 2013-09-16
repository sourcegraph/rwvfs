package rwvfs

import (
	"bytes"
	"fmt"
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
		testWrite(t, test.fs, test.path)
	}
}

func testWrite(t *testing.T, fs FileSystem, path string) {
	label := fmt.Sprintf("%T", fs)

	w, err := fs.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC|os.O_EXCL)
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
