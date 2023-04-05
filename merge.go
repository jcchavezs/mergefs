package mergefs

import (
	"errors"
	"io/fs"
	"os"
)

// Merge filesystems
func Merge(filesystems ...fs.FS) fs.FS {
	if len(filesystems) == 0 {
		return mergedFS{}
	}

	if len(filesystems) == 1 {
		return filesystems[0]
	}

	mfs := mergedFS{filesystems: filesystems}
	var (
		supportsGlob     = true
		supportsReadFile = true
	)

	for _, wfs := range filesystems {
		if _, ok := wfs.(globber); !ok {
			supportsGlob = false
			break
		}
	}

	for _, wfs := range filesystems {
		if _, ok := wfs.(fileReader); !ok {
			supportsReadFile = false
			break
		}
	}

	switch {
	case supportsGlob && supportsReadFile:
		return struct {
			mergedFS
			globber
			fileReader
		}{
			mfs,
			&globFS{mfs},
			&readFileFS{mfs},
		}
	case supportsGlob && !supportsReadFile:
		return struct {
			mergedFS
			globber
		}{
			mfs,
			&globFS{mfs},
		}
	case !supportsGlob && supportsReadFile:
		return struct {
			mergedFS
			fileReader
		}{
			mfs,
			&readFileFS{mfs},
		}
	case !supportsGlob && !supportsReadFile:
		return mfs
	}

	return mfs
}

// mergedFS combines filesystems. Each filesystem can serve different paths. The first FS takes precedence
type mergedFS struct {
	filesystems []fs.FS
}

// Open opens the named file.
func (mfs mergedFS) Open(name string) (fs.File, error) {
	for _, fs := range mfs.filesystems {
		file, err := fs.Open(name)
		if err == nil { // TODO should we return early when it's not an os.ErrNotExist? Should we offer options to decide this behaviour?
			return file, nil
		}
	}
	return nil, os.ErrNotExist
}

// ReadDir reads from the directory, and produces a DirEntry array of different
// directories.
//
// It iterates through all different filesystems that exist in the mfs MergeFS
// filesystem slice and it identifies overlapping directories that exist in different
// filesystems
func (mfs mergedFS) ReadDir(name string) ([]fs.DirEntry, error) {
	dirsMap := make(map[string]fs.DirEntry)
	notExistCount := 0
	for _, filesystem := range mfs.filesystems {
		dir, err := fs.ReadDir(filesystem, name)
		if err != nil {
			if errors.Is(err, fs.ErrNotExist) {
				notExistCount++
				continue
			}
			return nil, err
		}
		for _, v := range dir {
			if _, ok := dirsMap[v.Name()]; !ok {
				dirsMap[v.Name()] = v
			}
		}
		continue
	}
	if len(mfs.filesystems) == notExistCount {
		return nil, fs.ErrNotExist
	}
	dirs := make([]fs.DirEntry, 0, len(dirsMap))

	for _, value := range dirsMap {
		dirs = append(dirs, value)
	}

	return dirs, nil
}
