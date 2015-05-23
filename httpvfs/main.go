package main

import (
	"flag"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"sourcegraph.com/sourcegraph/rwvfs"
)

var (
	httpAddr   = flag.String("http", ":7070", "HTTP listen address")
	storageDir = flag.String("dir", ".", "directory to serve via HTTP (caution: HTTP clients can create/edit/delete/view any files in this dir)")
)

func main() {
	log.SetFlags(0)
	flag.Parse()

	dir, err := filepath.Abs(*storageDir)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Serving %s on %s", dir, *httpAddr)

	http.Handle("/", rwvfs.HTTPHandler(rwvfs.OS(dir), os.Stderr))
	log.Fatal(http.ListenAndServe(*httpAddr, nil))
}
