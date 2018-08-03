package building

import (
	"fmt"
	"path/filepath"
)

func glob(dir string, files []string, fail bool) ([]string, error) {
	var paths []string
	for _, f := range files {
		matches, err := filepath.Glob(filepath.Join(dir, f))
		if err != nil {
			return nil, err
		}
		if matches != nil {
			for _, match := range matches {
				if dir != "" {
					match, err = filepath.Rel(dir, match)
					if err != nil {
						return nil, err
					}
				}
				paths = append(paths, match)
			}
		} else if fail {
			return nil, fmt.Errorf("file %q does not exist", f)
		}
	}
	if fail && len(paths) == 0 {
		return nil, fmt.Errorf("no source files")
	}
	return paths, nil
}
