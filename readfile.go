package mergefs

import (
	"errors"
	"os"
)

type fileReader interface {
	ReadFile(name string) ([]byte, error)
}

type readFileFS struct {
	mergedFS
}

func (rfs *readFileFS) ReadFile(name string) (content []byte, err error) {
	for _, fs := range rfs.mergedFS.filesystems {
		frfs := fs.(fileReader)
		content, err = frfs.ReadFile(name)
		if err == nil {
			return content, nil
		} else if !errors.Is(err, os.ErrNotExist) {
			if errors.Unwrap(err).Error() != "invalid name" {
				// invalid name is returned when the filename
				// is an absolute path, which is only supported by
				// the OS filesystem
				return nil, err
			}
		}
	}

	return nil, err
}
