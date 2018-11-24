package building

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"
)

func makeFileset(dir, includes, excludes string) fileset {
	return fileset{
		dir:      dir,
		includes: strings.Split(includes, ","),
		excludes: strings.Split(excludes, ","),
	}
}

type fileset struct {
	dir      string
	includes []string
	excludes []string
}

func resolve(filesets []fileset, fail bool) ([]fileset, error) {
	var fs []fileset
	for _, f := range filesets {
		fss, err := f.glob(fail)
		if err != nil {
			return nil, err
		}
		fs = append(fs, fss)
	}
	return fs, nil
}

func (f fileset) glob(fail bool) (fileset, error) {
	var paths []string
	for _, include := range f.includes {
		if f.dir != "" && !filepath.IsAbs(include) {
			include = filepath.Join(f.dir, include)
		}
		include = filepath.Clean(include)
		matches, err := filepath.Glob(include)
		if err != nil {
			return f, err
		}
		if matches != nil {
			for _, match := range matches {
				if f.dir != "" {
					match, err = filepath.Rel(f.dir, match)
					if err != nil {
						return f, err
					}
				}
				paths = append(paths, filepath.ToSlash(match))
			}
		} else if fail {
			return f, fmt.Errorf("file %q does not exist", include)
		}
	}
	if fail && len(paths) == 0 {
		return f, fmt.Errorf("no source files")
	}
	f.includes = paths
	return f, nil
}

func (f fileset) walk(fn func(path, rel string, info os.FileInfo, err error) error) error {
	paths := make(map[string]bool)
	for _, include := range f.includes {
		if f.dir != "" && !filepath.IsAbs(include) {
			include = filepath.Join(f.dir, include)
		}
		if err := filepath.Walk(include, func(p string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			rel := p
			if f.dir != "" {
				rel, err = filepath.Rel(f.dir, rel)
				if err != nil {
					return err
				}
			} else {
				rel, err = filepath.Rel(include, rel)
				if err != nil {
					return err
				}
				rel = filepath.Join(filepath.Base(include), rel)
			}
			rel = filepath.ToSlash(rel)
			if paths[rel] {
				return nil
			}
			paths[rel] = true
			for _, exclude := range f.excludes {
				skip, err := path.Match(exclude, rel)
				if err != nil {
					return err
				}
				if skip {
					b.Debugf("excluded %q", p)
					return filepath.SkipDir
				}
			}
			return fn(p, rel, info, err)
		}); err != nil {
			return err
		}
	}
	return nil
}
