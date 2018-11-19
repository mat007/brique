package building

import (
	"archive/tar"
	"compress/gzip"
	"debug/elf"
	"io"
	"os"
	"path/filepath"
)

func (b *B) Tar(args ...string) compression {
	return makeCompression(Tar{}, args)
}

type Tar struct{}

func (t Tar) Write(w io.Writer, level int, dst string, srcs []fileset) error {
	if dst != "-" {
		f, err := os.Create(dst)
		if err != nil {
			return err
		}
		defer b.Close(f)
		w = f
		ext := filepath.Ext(dst)
		if ext == ".gz" || ext == ".tgz" {
			gz, err := gzip.NewWriterLevel(f, level)
			if err != nil {
				return err
			}
			defer b.Close(gz)
			w = gz
		}
	}
	tw := tar.NewWriter(w)
	defer b.Close(tw)
	return walk(srcs, func(path, rel string, info os.FileInfo) error {
		return writeTar(tw, path, rel, info)
	})
}

func writeTar(tw *tar.Writer, path, rel string, info os.FileInfo) error {
	hdr, err := tar.FileInfoHeader(info, info.Name())
	if err != nil {
		return err
	}
	hdr.Name = filepath.ToSlash(rel)
	if hdr.Mode%2 == 0 && isExecutable(path) {
		hdr.Mode++
		b.Debugln("fixed execute permissions for", hdr.Name)
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
	defer b.Close(f)
	_, err = io.Copy(tw, f)
	return err
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
		defer b.Close(f)
		r = f
		ext := filepath.Ext(src)
		if ext == ".gz" || ext == ".tgz" {
			gz, err := gzip.NewReader(r)
			if err != nil {
				return err
			}
			defer b.Close(gz)
			r = gz
		}
	}
	tr := tar.NewReader(r)
	for {
		hdr, err := tr.Next()
		if err == io.EOF {
			break // End of compression
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
		defer b.Close(f)
		_, err = io.Copy(f, tr)
		if err != nil {
			return err
		}
	}
	return nil
}
