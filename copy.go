package building

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
)

// Copy wraps a files and folders copy operation.
type Copy struct {
	destination string
	filesets    []fileset
}

// Copy handles copying files and folders.
func (b *B) Copy(destination string, paths ...string) Copy {
	c := Copy{
		destination: destination,
	}
	if len(paths) > 0 {
		c.filesets = append(c.filesets, fileset{
			includes: paths,
		})
		c.Run()
	}
	return c
}

// WithFiles adds files and folders to copy.
func (c Copy) WithFiles(paths ...string) Copy {
	c.filesets = append(c.filesets, fileset{
		includes: paths,
	})
	return c
}

// WithFileset adds a fileset to copy.
func (c Copy) WithFileset(dir, includes, excludes string) Copy {
	c.filesets = append(c.filesets, makeFileset(dir, includes, excludes))
	return c
}

// Run performs the copy.
func (c Copy) Run() {
	if err := c.run(); err != nil {
		b.Fatalln(err)
	}
}

func (c Copy) run() error {
	filesets, err := resolve(c.filesets, true)
	if err != nil {
		return err
	}
	if len(filesets) == 0 {
		return nil
	}
	toFile := true
	info, err := os.Stat(c.destination)
	if os.IsNotExist(err) {
		if c.destination[len(c.destination)-1] == '/' || len(filesets) > 1 || len(filesets[0].includes) > 1 {
			toFile = false
		}
	} else if err != nil {
		return err
	} else if info.IsDir() {
		toFile = false
	} else if len(filesets) > 1 || len(filesets[0].includes) > 1 {
		return fmt.Errorf("only one source file allowed when destination is a file")
	}
	for _, fs := range filesets {
		return fs.walk(func(path, rel string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			dest := filepath.Join(c.destination, rel)
			if info.IsDir() {
				toFile = false
				return os.MkdirAll(dest, info.Mode())
			}
			if toFile {
				dest = c.destination
			}
			dirInfo, err := os.Stat(filepath.Dir(path))
			if err != nil {
				return err
			}
			if err = os.MkdirAll(filepath.Dir(dest), dirInfo.Mode()); err != nil {
				return err
			}
			b.Debugf("copying file %q to %q\n", path, dest)
			return copyFile(path, dest, info.Mode())
		})
	}
	return nil
}

func copyFile(src, dst string, mode os.FileMode) error {
	if same, err := sameFile(src, dst); err != nil || same {
		return err
	}
	source, err := os.Open(src)
	if err != nil {
		return err
	}
	defer b.Close(source)
	destination, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer b.Close(destination)
	if _, err = io.Copy(destination, source); err != nil {
		return err
	}
	return os.Chmod(dst, mode)
}

func sameFile(src, dst string) (bool, error) {
	absSrc, err := filepath.Abs(src)
	if err != nil {
		return false, err
	}
	absDst, err := filepath.Abs(dst)
	if err != nil {
		return false, err
	}
	return absSrc == absDst, nil
}
