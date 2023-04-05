package mergefs

import (
	"embed"
	"io/fs"
	"os"
	"testing"
	"testing/fstest"

	"github.com/jcchavezs/mergefs/io"
	"github.com/stretchr/testify/require"
)

func TestReadfile(t *testing.T) {
	var aFS fs.ReadFileFS = fstest.MapFS{
		"testdata":              &fstest.MapFile{Mode: fs.ModeDir},
		"testdata/a":            &fstest.MapFile{Mode: fs.ModeDir},
		"testdata/a/y":          &fstest.MapFile{Mode: fs.ModeDir},
		"testdata/a/y/foo.conf": &fstest.MapFile{Data: []byte("bar")},
	}

	var bFS fs.ReadFileFS = fstest.MapFS{
		"testdata":              &fstest.MapFile{Mode: fs.ModeDir},
		"testdata/b":            &fstest.MapFile{Mode: fs.ModeDir},
		"testdata/b/y":          &fstest.MapFile{Mode: fs.ModeDir},
		"testdata/b/y/bar.conf": &fstest.MapFile{Data: []byte("bar")},
	}

	mfs := Merge(aFS, bFS)

	t.Run("file exists", func(t *testing.T) {
		content, err := fs.ReadFile(mfs, "testdata/a/y/foo.conf")
		require.NoError(t, err)
		require.Equal(t, "bar", string(content))

		rfmfs, ok := mfs.(fs.ReadFileFS)
		require.True(t, ok)
		content, err = rfmfs.ReadFile("testdata/a/y/foo.conf")
		require.NoError(t, err)
		require.Equal(t, "bar", string(content))
	})

	t.Run("file does not exist", func(t *testing.T) {
		_, err := fs.ReadFile(mfs, "testdata/a/y/baz.conf")
		require.Error(t, err)

		rfmfs, ok := mfs.(fs.ReadFileFS)
		require.True(t, ok)
		_, err = rfmfs.ReadFile("testdata/a/y/baz.conf")
		require.Error(t, err)
	})
}

//go:embed testdata
var testdataFS embed.FS

func TestAbsolutePath(t *testing.T) {
	mfs := Merge(testdataFS, io.OSFS)

	f, err := os.CreateTemp(t.TempDir(), "fizz.conf")
	require.NoError(t, err)
	defer f.Close()

	require.Equal(t, string(f.Name()[0]), "/")

	_, err = fs.ReadFile(mfs, f.Name())
	require.NoError(t, err)

	rfmfs, ok := mfs.(fs.ReadFileFS)
	require.True(t, ok)

	_, err = rfmfs.ReadFile(f.Name())
	require.NoError(t, err)
}
