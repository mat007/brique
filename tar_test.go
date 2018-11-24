package building

import (
	"fmt"
	"path/filepath"
	"testing"

	"gotest.tools/assert"
	"gotest.tools/fs"
)

func TestTarTree(t *testing.T) {
	rootDirectory := fs.NewDir(t, "root",
		fs.WithDir("source",
			fs.WithFile("foo.txt", "foo"),
			fs.WithDir("bar",
				fs.WithFile("bar.txt", "bar"))))
	defer rootDirectory.Remove()

	dst := filepath.Join(rootDirectory.Path(), "destination", "dst.tar")
	err := compress(Tar{}, nil, -1, dst,
		fileset{dir: rootDirectory.Path() + "/source", includes: []string{"foo.txt"}},
		fileset{dir: rootDirectory.Path(), includes: []string{"source/bar"}})
	assert.NilError(t, err)
	err = untarFiles(dst, filepath.Join(rootDirectory.Path(), "destination"))
	assert.NilError(t, err)

	expected := fs.Expected(t,
		fs.WithDir("source",
			fs.WithFile("foo.txt", "foo"),
			fs.WithDir("bar",
				fs.WithFile("bar.txt", "bar"))),
		fs.WithDir("destination",
			fs.WithFile("dst.tar", "", fs.MatchAnyFileContent),
			fs.WithFile("foo.txt", "foo"),
			fs.WithDir("source",
				fs.WithDir("bar",
					fs.WithFile("bar.txt", "bar")))))
	assert.Assert(t, fs.Equal(rootDirectory.Path(), expected))
}

func TestTarTreeWithGlob(t *testing.T) {
	rootDirectory := fs.NewDir(t, "root",
		fs.WithDir("source",
			fs.WithFile("foo.txt", "foo"),
			fs.WithDir("bar",
				fs.WithFile("bar.txt", "bar"))))
	defer rootDirectory.Remove()

	dst := filepath.Join(rootDirectory.Path(), "destination", "dst.tar")
	err := compress(Tar{}, nil, -1, dst, fileset{dir: rootDirectory.Path(), includes: []string{"source/*"}})
	assert.NilError(t, err)
	err = untarFiles(dst, filepath.Join(rootDirectory.Path(), "destination"))
	assert.NilError(t, err)

	expected := fs.Expected(t,
		fs.WithDir("source",
			fs.WithFile("foo.txt", "foo"),
			fs.WithDir("bar",
				fs.WithFile("bar.txt", "bar"))),
		fs.WithDir("destination",
			fs.WithFile("dst.tar", "", fs.MatchAnyFileContent),
			fs.WithDir("source",
				fs.WithFile("foo.txt", "foo"),
				fs.WithDir("bar",
					fs.WithFile("bar.txt", "bar")))))
	assert.Assert(t, fs.Equal(rootDirectory.Path(), expected))
}

func TestTarTreeWithEmptyGlob(t *testing.T) {
	rootDirectory := fs.NewDir(t, "root",
		fs.WithDir("source",
			fs.WithFile("foo.txt", "foo"),
			fs.WithDir("bar",
				fs.WithFile("bar.txt", "bar"))))
	defer rootDirectory.Remove()

	dst := filepath.Join(rootDirectory.Path(), "destination", "dst.tar")
	src := "source/non-existing*"
	err := compress(Tar{}, nil, -1, dst, fileset{dir: rootDirectory.Path(), includes: []string{src}})
	assert.Error(t, err, fmt.Sprintf("file %q does not exist", filepath.Join(rootDirectory.Path(), src)))
}

func TestGzipTarTree(t *testing.T) {
	rootDirectory := fs.NewDir(t, "root",
		fs.WithDir("source",
			fs.WithFile("foo.txt", "foo"),
			fs.WithDir("bar",
				fs.WithFile("bar.txt", "bar"))))
	defer rootDirectory.Remove()

	dst := filepath.Join(rootDirectory.Path(), "destination", "dst.tar.gz")
	err := compress(Tar{}, nil, -1, dst,
		fileset{dir: rootDirectory.Path(), includes: []string{"source/foo.txt", "source/bar"}})
	assert.NilError(t, err)
	err = untarFiles(dst, filepath.Join(rootDirectory.Path(), "destination"))
	assert.NilError(t, err)

	expected := fs.Expected(t,
		fs.WithDir("source",
			fs.WithFile("foo.txt", "foo"),
			fs.WithDir("bar",
				fs.WithFile("bar.txt", "bar"))),
		fs.WithDir("destination",
			fs.WithFile("dst.tar.gz", "", fs.MatchAnyFileContent),
			fs.WithDir("source",
				fs.WithFile("foo.txt", "foo"),
				fs.WithDir("bar",
					fs.WithFile("bar.txt", "bar")))))
	assert.Assert(t, fs.Equal(rootDirectory.Path(), expected))
}

func TestTarTreeWithExcludes(t *testing.T) {
	rootDirectory := fs.NewDir(t, "root",
		fs.WithDir("source",
			fs.WithFile("foo.txt", "foo"),
			fs.WithDir("bar",
				fs.WithFile("bar.txt", "bar"))))
	defer rootDirectory.Remove()

	dst := filepath.Join(rootDirectory.Path(), "destination", "dst.tar")
	err := compress(Tar{}, nil, -1, dst, fileset{dir: rootDirectory.Path(), includes: []string{"*"}, excludes: []string{"*/foo*"}})
	assert.NilError(t, err)
	err = untarFiles(dst, filepath.Join(rootDirectory.Path(), "destination"))
	assert.NilError(t, err)

	expected := fs.Expected(t,
		fs.WithDir("source",
			fs.WithFile("foo.txt", "foo"),
			fs.WithDir("bar",
				fs.WithFile("bar.txt", "bar"))),
		fs.WithDir("destination",
			fs.WithFile("dst.tar", "", fs.MatchAnyFileContent),
			fs.WithDir("source",
				fs.WithDir("bar",
					fs.WithFile("bar.txt", "bar")))))
	assert.Assert(t, fs.Equal(rootDirectory.Path(), expected))
}
