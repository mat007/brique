package building

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
)

type archive interface {
	Write(w io.Writer, level int, dst string, srcs []fileset) error
}

type compression struct {
	output   io.Writer
	filesets []fileset
	archive  archive
	level    int
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

// WithFiles adds files to compress.
func (t compression) WithFiles(paths ...string) compression {
	t.filesets = append(t.filesets, fileset{
		includes: paths,
	})
	return t
}

// WithFileset adds a fileset to compress.
func (t compression) WithFileset(dir, includes, excludes string) compression {
	t.filesets = append(t.filesets, makeFileset(dir, includes, excludes))
	return t
}

// WithLevel sets the compression level from 0 (no compression) to 9 (best compression).
// -1 can be used for default compression level.
func (t compression) WithLevel(level int) compression {
	if level < -1 && level > 9 {
		b.Fatalln("invalid compression level", level)
	}
	t.level = level
	return t
}

func (t compression) Run(dst string, args ...string) {
	if len(args) > 0 {
		t.filesets = append(t.filesets, fileset{
			includes: args,
		})
	}
	if t.output == nil {
		t.output = os.Stdout
	}
	if err := compress(t.archive, t.output, t.level, dst, t.filesets...); err != nil {
		b.Fatalln(err)
	}
}

func compress(a archive, w io.Writer, level int, dst string, srcs ...fileset) error {
	fs, err := resolve(srcs, true)
	if err != nil {
		return err
	}
	if len(fs) == 0 {
		return fmt.Errorf("needs at least one file")
	}
	if dst != "-" {
		if err := os.MkdirAll(filepath.Dir(dst), 0755); err != nil {
			return err
		}
	}
	return a.Write(w, level, dst, fs)
}

func walk(fs []fileset, fn func(path, rel string, info os.FileInfo) error) error {
	for _, f := range fs {
		err := f.walk(func(path, rel string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			b.Debugln("compressing", rel)
			return fn(path, rel, info)
		})
		if err != nil {
			return err
		}
	}
	return nil
}
