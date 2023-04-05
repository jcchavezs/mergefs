package mergefs

type globber interface {
	Glob(pattern string) (matches []string, err error)
}

type globFS struct {
	mergedFS
}

func (gfs *globFS) Glob(pattern string) (matches []string, err error) {
	var currentMatches []string
	var matchesMap = make(map[string]struct{})
	for _, fs := range gfs.mergedFS.filesystems {
		gfs := fs.(globber)
		currentMatches, err = gfs.Glob(pattern)
		if err != nil {
			return nil, err
		}

		for _, match := range currentMatches {
			if _, ok := matchesMap[match]; !ok {
				matchesMap[match] = struct{}{}
				matches = append(matches, match)
			}
		}
	}

	return
}
