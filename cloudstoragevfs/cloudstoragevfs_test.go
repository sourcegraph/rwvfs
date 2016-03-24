package cloudstoragevfs_test

import (
	"os"
	"testing"

	"golang.org/x/net/context"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/storage/v1"

	"sourcegraph.com/sourcegraph/rwvfs/cloudstoragevfs"
	"sourcegraph.com/sourcegraph/rwvfs/testutil"
)

func TestVFS(t *testing.T) {
	if os.Getenv("TRAVIS") == "true" {
		t.Skip("no credentials on Travis")
	}

	_, err := google.DefaultClient(context.Background(), storage.DevstorageFullControlScope)
	if err != nil {
		t.Skip(err)
	}

	fs, err := cloudstoragevfs.NewDefault("rwvfs-test")
	if err != nil {
		t.Error(err)
	}
	testutil.Write(t, fs)
	testutil.Glob(t, fs)
}
