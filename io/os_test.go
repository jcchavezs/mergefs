package io

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestOpen(t *testing.T) {
	f, err := os.CreateTemp(t.TempDir(), "fizz.conf")
	require.NoError(t, err)
	defer f.Close()

	_, err = OSFS.Open(f.Name())
	require.NoError(t, err)

	_, err = OSFS.Open("testdata/fuzz.conf")
	require.NoError(t, err)

	_, err = OSFS.Open("testdata/bar.conf")
	require.Error(t, err)
}

func TestReadFile(t *testing.T) {
	f, err := os.CreateTemp(t.TempDir(), "fizz.conf")
	require.NoError(t, err)
	defer f.Close()

	_, err = OSFS.ReadFile(f.Name())
	require.NoError(t, err)

	_, err = OSFS.ReadFile("testdata/fuzz.conf")
	require.NoError(t, err)

	_, err = OSFS.ReadFile("testdata/bar.conf")
	require.Error(t, err)
}

func TestReadDir(t *testing.T) {
	_, err := OSFS.ReadDir(t.TempDir())
	require.NoError(t, err)

	_, err = OSFS.ReadDir("testdata")
	require.NoError(t, err)

	_, err = OSFS.ReadDir("notestdata")
	require.Error(t, err)
}

func TestGlob(t *testing.T) {
	matches, err := OSFS.Glob("testdata/*.conf")
	require.NoError(t, err)
	require.Len(t, matches, 1)

	matches, err = OSFS.Glob("notestdata/*.conf")
	require.NoError(t, err)
	require.Len(t, matches, 0)
}
