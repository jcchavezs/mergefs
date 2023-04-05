package io

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
)

// Inspired in https://github.com/corazawaf/coraza/blob/v3/dev/internal/io/file.go

// OSFS implements fs.FS using methods on OS to read from the system.
// More context in: https://github.com/golang/go/issues/44279
type osFS struct{}

func (osFS) Open(name string) (fs.File, error) {
	if name[0] == '/' {
		return os.Open(name)
	}
	// TODO(jcchavezs): Shall we use ValidPath?
	absName, err := filepath.Abs(name)
	if err != nil {
		return nil, err
	}
	return os.Open(absName)
}

func (osFS) ReadFile(name string) ([]byte, error) {
	if name[0] == '/' {
		return os.ReadFile(name)
	}

	// TODO(jcchavezs): Shall we use ValidPath?
	absName, err := filepath.Abs(name)
	if err != nil {
		fmt.Println(absName)
		return nil, err
	}
	return os.ReadFile(name)
}

func (osFS) ReadDir(name string) ([]fs.DirEntry, error) {
	return os.ReadDir(name)
}

func (osFS) Glob(pattern string) ([]string, error) {
	return filepath.Glob(pattern)
}

var OSFS osFS
