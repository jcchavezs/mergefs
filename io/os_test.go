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

	_, err = OSFS.Open("./testdata/fuzz.conf")
	require.NoError(t, err)
}
