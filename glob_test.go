package mergefs

import (
	"io/fs"
	"testing"
	"testing/fstest"

	"github.com/stretchr/testify/require"
)

func TestGlob(t *testing.T) {
	var aFS fs.GlobFS = fstest.MapFS{
		"testdata":              &fstest.MapFile{Mode: fs.ModeDir},
		"testdata/a":            &fstest.MapFile{Mode: fs.ModeDir},
		"testdata/a/y":          &fstest.MapFile{Mode: fs.ModeDir},
		"testdata/a/y/foo.conf": &fstest.MapFile{Data: []byte("bar")},
	}

	var bFS fs.GlobFS = fstest.MapFS{
		"testdata":              &fstest.MapFile{Mode: fs.ModeDir},
		"testdata/b":            &fstest.MapFile{Mode: fs.ModeDir},
		"testdata/b/y":          &fstest.MapFile{Mode: fs.ModeDir},
		"testdata/b/y/bar.conf": &fstest.MapFile{Data: []byte("bar")},
	}

	mfs := Merge(aFS, bFS)

	matches, err := fs.Glob(mfs, "testdata/*/*/*.conf")
	require.NoError(t, err)
	require.Len(t, matches, 2)

	gmfs, ok := mfs.(fs.GlobFS)
	require.True(t, ok)
	matches, err = gmfs.Glob("testdata/*/*/*.conf")
	require.NoError(t, err)
	require.Len(t, matches, 2)
}
