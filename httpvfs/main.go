package main

import (
	"flag"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"sourcegraph.com/sourcegraph/rwvfs"
)

var (
	httpAddr   = flag.String("http", ":7070", "HTTP listen address")
	storageDir = flag.String("dir", ".", "directory to serve via HTTP (caution: HTTP clients can create/edit/delete/view any files in this dir)")
	verbose    = flag.Bool("v", false, "verbose output")
)

func main() {
	log.SetFlags(0)
	flag.Parse()

	dir, err := filepath.Abs(*storageDir)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Serving %s on %s", dir, *httpAddr)

	var logTo io.Writer
	if *verbose {
		logTo = os.Stderr
	} else {
		logTo = ioutil.Discard
	}

	http.Handle("/", rwvfs.HTTPHandler(rwvfs.OS(dir), logTo))
	log.Fatal(http.ListenAndServe(*httpAddr, nil))
}
