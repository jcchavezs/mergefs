package mergefs_test

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"testing"
	"testing/fstest"

	"github.com/jcchavezs/mergefs"
	"github.com/stretchr/testify/require"
)

func TestMerge(t *testing.T) {
	t.Run("empty merge", func(t *testing.T) {
		mfs := mergefs.Merge()

		if _, err := mfs.Open("a"); err == nil || !errors.Is(err, os.ErrNotExist) {
			t.Fatalf("expected error")
		}
	})

	t.Run("merge of one FS", func(t *testing.T) {
		sfs := fstest.MapFS{}
		mfs := mergefs.Merge(sfs)
		require.Equal(t, fmt.Sprintf("%v", sfs), fmt.Sprintf("%v", mfs))
	})

	t.Run("merge of two FS", func(t *testing.T) {
		a := fstest.MapFS{"a": &fstest.MapFile{Data: []byte("text")}}
		b := fstest.MapFS{"b": &fstest.MapFile{Data: []byte("text")}}
		mfs := mergefs.Merge(a, b)

		if _, err := mfs.Open("a"); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if _, err := mfs.Open("b"); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if _, err := mfs.Open("c"); err == nil {
			t.Fatalf("expected error")
		}
	})
}

var (
	aFS = fstest.MapFS{
		"a":            &fstest.MapFile{Mode: fs.ModeDir},
		"a/z":          &fstest.MapFile{Mode: fs.ModeDir},
		"a/z/bar.conf": &fstest.MapFile{Data: []byte("bar")},
	}
	bFS = os.DirFS(filepath.Join("testdata"))

	abFS = mergefs.Merge(aFS, bFS)
)

func TestMergedReadDir(t *testing.T) {
	var filePaths = []struct {
		path           string
		dirArrayLength int
		child          string
	}{
		// MapFS takes in account the current directory in addition to all included directories and produces a "" dir
		{"a", 1, "z"},
		{"a/z", 1, "bar.conf"},
		{"b", 1, "z"},
		{"b/z", 1, "foo.conf"},
		{"c", 0, ""},
	}

	for _, fp := range filePaths {
		t.Run(fp.path, func(t *testing.T) {
			dirs, err := fs.ReadDir(abFS, fp.path)
			if fp.dirArrayLength > 0 {
				require.NoError(t, err)
			}
			require.Len(t, dirs, fp.dirArrayLength)

			for i := 0; i < len(dirs); i++ {
				require.Equal(t, dirs[i].Name(), fp.child)
			}
		})
	}
}

func TestMergedOpen(t *testing.T) {
	data := make([]byte, 3)
	file, err := abFS.Open("a/z/bar.conf")
	require.NoError(t, err)

	_, err = file.Read(data)
	require.NoError(t, err)
	require.Equal(t, "bar", string(data))

	file, err = abFS.Open("b/z/foo.conf")
	require.NoError(t, err)

	_, err = file.Read(data)
	require.NoError(t, err)
	require.Equal(t, "foo", string(data))

	require.NoError(t, file.Close())
}
