package building

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
)

type fileset struct {
	dir   string
	files []string
}

type archive interface {
	Write(w io.Writer, level int, dst string, srcs []fileset) error
}

type compression struct {
	output  io.Writer
	files   []fileset
	archive archive
	level   int
}

func makeCompression(a archive, args []string) compression {
	c := compression{
		archive: a,
		level:   -1,
	}
	if len(args) > 0 {
		c.Run(args[0], args[1:]...)
	}
	return c
}

func (t compression) WithOutput(w io.Writer) compression {
	t.output = w
	return t
}

func (t compression) WithFiles(dir string, files ...string) compression {
	if len(files) > 0 {
		t.files = append(t.files, fileset{
			dir:   dir,
			files: files,
		})
	}
	return t
}

// WithLevel sets the compression level from 0 (no compression) to 9 (best compression).
// -1 can be used for default compression level.
func (t compression) WithLevel(level int) compression {
	if level < -1 && level > 9 {
		Fatalln("invalid compression level", level)
	}
	t.level = level
	return t
}

func (t compression) Run(dst string, args ...string) {
	if len(args) > 0 {
		t.files = append(t.files, fileset{
			files: args,
		})
	}
	if t.output == nil {
		t.output = os.Stdout
	}
	if err := compress(t.archive, t.output, t.level, dst, t.files...); err != nil {
		Fatalln(err)
	}
}

func compress(a archive, w io.Writer, level int, dst string, srcs ...fileset) error {
	files := []fileset{}
	for _, f := range srcs {
		matches, err := glob(f.dir, f.files, true)
		if err != nil {
			return err
		}
		files = append(files, fileset{
			dir:   f.dir,
			files: matches,
		})
	}
	if len(files) == 0 {
		return fmt.Errorf("needs at least one file")
	}
	if dst != "-" {
		if err := os.MkdirAll(filepath.Dir(dst), 0755); err != nil {
			return err
		}
	}
	return a.Write(w, level, dst, files)
}

func walk(srcs []fileset, f func(path, rel string, info os.FileInfo) error) error {
	for _, src := range srcs {
		dir := src.dir
		for _, file := range src.files {
			err := filepath.Walk(filepath.Join(dir, file), func(path string, info os.FileInfo, err error) error {
				if err != nil {
					return err
				}
				rel, err := filepath.Rel(dir, path)
				if err != nil {
					return err
				}
				return f(path, rel, info)
			})
			if err != nil {
				return err
			}
		}
	}
	return nil
}
