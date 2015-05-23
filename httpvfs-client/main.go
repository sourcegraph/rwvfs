package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net/url"
	"os"
	"path"
	"time"

	"strings"

	"sourcegraph.com/sourcegraph/rwvfs"
)

var (
	urlStr = flag.String("url", "http://localhost:7070/", "URL to HTTP VFS (typically served by the `httpvfs` program)")
)

func main() {
	log.SetFlags(0)
	flag.Parse()

	if flag.NArg() != 2 {
		log.Fatal("error: usage: httpvfs-client [opts] <cat|ls|put|rm> <path>")
	}
	op := flag.Arg(0)
	path := path.Clean(flag.Arg(1))

	url, err := url.Parse(*urlStr)
	if err != nil {
		log.Fatal(err)
	}

	fs := rwvfs.HTTP(url, nil)

	switch strings.ToLower(op) {
	case "cat":
		f, err := fs.Open(path)
		if err != nil {
			log.Fatal(err)
		}
		defer func() {
			if err := f.Close(); err != nil {
				log.Fatal(err)
			}
		}()
		if _, err := io.Copy(os.Stdout, f); err != nil {
			log.Fatal(err)
		}

	case "ls":
		fis, err := fs.ReadDir(path)
		if err != nil {
			log.Fatal(err)
		}

		var longestNameLen int
		for _, fi := range fis {
			if len(fi.Name()) > longestNameLen {
				longestNameLen = len(fi.Name())
			}
		}
		longestNameLen++ // account for "/" suffix on dirs

		for _, fi := range fis {
			name := fi.Name()
			if fi.IsDir() {
				name += "/"
			}
			mtime := fi.ModTime().Round(time.Second)
			fmt.Printf("%-*s   %s   %d\n", longestNameLen, name, mtime, fi.Size())
		}

	case "put":
		f, err := fs.Create(path)
		if err != nil {
			log.Fatal(err)
		}
		defer func() {
			if err := f.Close(); err != nil {
				log.Fatal(err)
			}
		}()
		log.Println("(reading file data on stdin...)")
		if _, err := io.Copy(f, os.Stdin); err != nil {
			log.Fatal(err)
		}

	case "rm":
		if err := fs.Remove(path); err != nil {
			log.Fatal(err)
		}

	default:
		log.Fatal("error: invalid op (see -h)")
	}
}
