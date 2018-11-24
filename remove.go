package building

import (
	"io"
	"os"
	"path/filepath"
)

// Remove wraps a files and folders remove operation.
type Remove struct {
	filesets  []fileset
	keepGoing bool
}

// Remove handles files and folders deletion.
func (b *B) Remove(paths ...string) Remove {
	r := Remove{}
	if len(paths) > 0 {
		r.filesets = append(r.filesets, fileset{
			includes: paths,
		})
		r.Run()
	}
	return r
}

// WithFiles adds a list of files and folders for deletion.
func (r Remove) WithFiles(paths ...string) Remove {
	r.filesets = append(r.filesets, fileset{
		includes: paths,
	})
	return r
}

// WithFileset adds a fileset for deletion.
func (r Remove) WithFileset(dir, includes, excludes string) Remove {
	r.filesets = append(r.filesets, makeFileset(dir, includes, excludes))
	return r
}

// Run performs the deletion.
func (r Remove) Run(paths ...string) {
	r.filesets = append(r.filesets, fileset{
		includes: paths,
	})
	if err := r.run(); err != nil {
		b.Fatalln(err)
	}
}

func (r Remove) run() error {
	filesets, err := resolve(r.filesets, false)
	if err != nil {
		return err
	}
	for _, f := range filesets {
		if len(f.excludes) == 0 {
			if err := removeWithoutExcludes(f); err != nil {
				return err
			}
		} else if err := removeWithExcludes(f); err != nil {
			return err
		}
	}
	return nil
}

func removeWithoutExcludes(f fileset) error {
	for _, include := range f.includes {
		if f.dir != "" && !filepath.IsAbs(include) {
			include = filepath.Join(f.dir, include)
		}
		err := os.RemoveAll(include)
		if err != nil {
			return err
		}
	}
	return nil
}

func removeWithExcludes(f fileset) error {
	var paths []string
	if err := f.walk(func(path, rel string, info os.FileInfo, err error) error {
		paths = append(paths, path)
		return err
	}); err != nil {
		return err
	}
	for i := len(paths); i > 0; i-- {
		path := paths[i-1]
		info, err := os.Stat(path)
		if os.IsNotExist(err) {
			b.Debugf("skipped non existing %q", path)
			continue
		}
		if err != nil {
			return err
		}
		if info.IsDir() {
			empty, err := isEmptyDir(path)
			if err != nil {
				return err
			}
			if !empty {
				b.Debugf("skipped non empty folder %q", path)
				continue
			}
		}
		if err := os.Remove(path); err != nil {
			return err
		}
		b.Debugf("removed %q", path)
	}
	return nil
}

func isEmptyDir(dir string) (bool, error) {
	f, err := os.Open(dir)
	if err != nil {
		return false, err
	}
	defer f.Close()
	_, err = f.Readdir(1)
	return err == io.EOF, nil
}
