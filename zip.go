package building

import (
	"archive/zip"
	"compress/flate"
	"errors"
	"io"
	"os"
	"path/filepath"
)

func (b *B) Zip(args ...string) compression {
	return makeCompression(Zip{}, args)
}

type Zip struct{}

func (z Zip) Name() string {
	return "zip"
}

func (z Zip) Write(w io.Writer, level int, dst string, srcs []fileset) error {
	if dst != "-" {
		f, err := os.Create(dst)
		if err != nil {
			return err
		}
		defer Close(f)
		w = f
	}
	zw := zip.NewWriter(w)
	zw.RegisterCompressor(zip.Deflate, func(w io.Writer) (io.WriteCloser, error) {
		return flate.NewWriter(w, level)
	})
	defer Close(zw)
	return walk(srcs, func(path, rel string, info os.FileInfo) error {
		if info.IsDir() {
			return nil
		}
		w, err := zw.Create(rel)
		if err != nil {
			return err
		}
		f, err := os.Open(path)
		if err != nil {
			return err
		}
		defer Close(f)
		_, err = io.Copy(w, f)
		return err
	})
}

func unzipFiles(src, dst string) error {
	var zr *zip.Reader
	if src == "-" {
		// $$$$ MAT to do
		return errors.New("reading from buffer not supported")
	} else {
		r, err := zip.OpenReader(src)
		if err != nil {
			return err
		}
		defer Close(r)
		zr = &r.Reader
	}
	for _, file := range zr.File {
		path := filepath.Join(dst, file.Name)
		info := file.FileInfo()
		if info.IsDir() {
			// $$$$ MAT pretty sure zip has only files
			if err := os.MkdirAll(path, info.Mode()); err != nil {
				return err
			}
			continue
		}
		dir := filepath.Dir(path)
		if err := os.MkdirAll(dir, 0755); err != nil {
			return err
		}
		f, err := os.OpenFile(path, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, info.Mode())
		if err != nil {
			return err
		}
		defer Close(f)
		r, err := file.Open()
		if err != nil {
			return err
		}
		defer Close(r)
		_, err = io.Copy(f, r)
		if err != nil {
			return err
		}
	}
	return nil
}
