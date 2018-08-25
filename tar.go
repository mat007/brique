package building

import (
	"archive/tar"
	"compress/gzip"
	"debug/elf"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

type tarr struct {
	output io.Writer
	files  []tarf
}

type tarf struct {
	dir   string
	files []string
}

func (b *B) Tar(args ...string) tarr {
	t := tarr{}
	if len(args) > 0 {
		t.Run(args...)
	}
	return t
}

func (t tarr) WithOutput(w io.Writer) tarr {
	t.output = w
	return t
}

func (t tarr) WithFiles(dir string, files ...string) tarr {
	if len(files) > 0 {
		t.files = append(t.files, tarf{
			dir:   dir,
			files: files,
		})
	}
	return t
}

func (t tarr) Run(args ...string) {
	if len(args) == 0 {
		Fatal("tar failed: needs at least one argument")
	}
	if len(args[1:]) > 0 {
		t.files = append(t.files, tarf{
			files: args[1:],
		})
	}
	if t.output == nil {
		t.output = os.Stdout
	}
	if err := tarFiles(t.output, args[0], t.files...); err != nil {
		Fatalln("tar failed:", err)
	}
}

func tarFiles(w io.Writer, dst string, srcs ...tarf) error {
	files := []tarf{}
	for _, f := range srcs {
		matches, err := glob(f.dir, f.files, true)
		if err != nil {
			return err
		}
		files = append(files, tarf{
			dir:   f.dir,
			files: matches,
		})
	}
	if len(files) == 0 {
		return fmt.Errorf("needs at least one file")
	}
	Debugln("tar", dst, files)
	if dst != "-" {
		if err := os.MkdirAll(filepath.Dir(dst), 0755); err != nil {
			return err
		}
		f, err := os.Create(dst)
		if err != nil {
			return err
		}
		defer f.Close()
		w = f
		ext := filepath.Ext(dst)
		if ext == ".gz" || ext == ".tgz" {
			gz := gzip.NewWriter(w)
			defer gz.Close()
			w = gz
		}
	}
	return writeTarFiles(w, files)
}

func writeTarFiles(w io.Writer, srcs []tarf) error {
	tw := tar.NewWriter(w)
	defer tw.Close()
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
				hdr, err := tar.FileInfoHeader(info, info.Name())
				if err != nil {
					return err
				}
				hdr.Name = filepath.ToSlash(rel)
				if hdr.Mode%2 == 0 && isExecutable(path) {
					hdr.Mode++
					Debugln("fixed execute permissions for", hdr.Name)
				}
				if err := tw.WriteHeader(hdr); err != nil {
					return err
				}
				if info.IsDir() {
					return nil
				}
				f, err := os.Open(path)
				if err != nil {
					return err
				}
				defer f.Close()
				_, err = io.Copy(tw, f)
				return err
			})
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func isExecutable(path string) bool {
	if e, err := elf.Open(path); err == nil && e.Type == elf.ET_EXEC {
		return true
	}
	if f, err := os.Open(path); err == nil {
		buf := make([]byte, 2)
		_, err = f.Read(buf)
		if err == nil && buf[0] == '#' && buf[1] == '!' {
			return true
		}
	}
	return false
}

func untarFiles(src, dst string) error {
	var r io.Reader
	if src == "-" {
		r = os.Stdin
	} else {
		f, err := os.Open(src)
		if err != nil {
			return err
		}
		defer f.Close()
		if err := os.MkdirAll(dst, 0755); err != nil {
			return err
		}
		r = f
		ext := filepath.Ext(src)
		if ext == ".gz" || ext == ".tgz" {
			gz, err := gzip.NewReader(r)
			if err != nil {
				return err
			}
			defer gz.Close()
			r = gz
		}
	}
	tr := tar.NewReader(r)
	for {
		hdr, err := tr.Next()
		if err == io.EOF {
			break // End of archive
		}
		if err != nil {
			return err
		}
		path := filepath.Join(dst, hdr.Name)
		info := hdr.FileInfo()
		if info.IsDir() {
			if err = os.MkdirAll(path, info.Mode()); err != nil {
				return err
			}
			continue
		}
		dir := filepath.Dir(path)
		if err = os.MkdirAll(dir, 0755); err != nil {
			return err
		}
		f, err := os.OpenFile(path, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, info.Mode())
		if err != nil {
			return err
		}
		defer f.Close()
		_, err = io.Copy(f, tr)
		if err != nil {
			return err
		}
	}
	return nil
}
