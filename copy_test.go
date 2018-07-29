package building

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"gotest.tools/assert"
	"gotest.tools/fs"
)

func TestCopyFileToNonExistingFile(t *testing.T) {
	rootDirectory := fs.NewDir(t, "root",
		fs.WithFile("foo.txt", "foo"))
	defer rootDirectory.Remove()

	err := copy(filepath.Join(rootDirectory.Path(), "bar.txt"),
		[]string{filepath.Join(rootDirectory.Path(), "foo.txt")})
	assert.NilError(t, err)

	expected := fs.Expected(t,
		fs.WithFile("foo.txt", "foo"),
		fs.WithFile("bar.txt", "foo"))
	assert.Assert(t, fs.Equal(rootDirectory.Path(), expected))
}

func TestCopyFileToExistingFile(t *testing.T) {
	rootDirectory := fs.NewDir(t, "root",
		fs.WithFile("foo.txt", "foo"),
		fs.WithFile("bar.txt", "bar"))
	defer rootDirectory.Remove()

	err := copy(filepath.Join(rootDirectory.Path(), "bar.txt"),
		[]string{filepath.Join(rootDirectory.Path(), "foo.txt")})
	assert.NilError(t, err)

	expected := fs.Expected(t,
		fs.WithFile("foo.txt", "foo"),
		fs.WithFile("bar.txt", "foo"))
	assert.Assert(t, fs.Equal(rootDirectory.Path(), expected))
}

func TestCopyFileToNonExistingDir(t *testing.T) {
	rootDirectory := fs.NewDir(t, "root",
		fs.WithFile("foo.txt", "foo"))
	defer rootDirectory.Remove()

	err := copy(filepath.Join(rootDirectory.Path(), "destination")+"/",
		[]string{filepath.Join(rootDirectory.Path(), "foo.txt")})
	assert.NilError(t, err)

	info, err := os.Stat(rootDirectory.Path())
	assert.NilError(t, err)
	expected := fs.Expected(t,
		fs.WithFile("foo.txt", "foo"),
		fs.WithDir("destination",
			fs.WithMode(info.Mode()),
			fs.WithFile("foo.txt", "foo")))
	assert.Assert(t, fs.Equal(rootDirectory.Path(), expected))
}

func TestCopyFileToExistingDir(t *testing.T) {
	rootDirectory := fs.NewDir(t, "root",
		fs.WithFile("foo.txt", "foo"),
		fs.WithDir("destination"))
	defer rootDirectory.Remove()

	err := copy(filepath.Join(rootDirectory.Path(), "destination"),
		[]string{filepath.Join(rootDirectory.Path(), "foo.txt")})
	assert.NilError(t, err)

	expected := fs.Expected(t,
		fs.WithFile("foo.txt", "foo"),
		fs.WithDir("destination",
			fs.WithFile("foo.txt", "foo")))
	assert.Assert(t, fs.Equal(rootDirectory.Path(), expected))
}

func TestCopyMultipleFilesToFile(t *testing.T) {
	rootDirectory := fs.NewDir(t, "root",
		fs.WithFile("foo.txt", "foo"),
		fs.WithFile("bar.txt", "bar"))
	defer rootDirectory.Remove()

	err := copy(filepath.Join(rootDirectory.Path(), "bar.txt"),
		[]string{
			filepath.Join(rootDirectory.Path(), "foo.txt"),
			filepath.Join(rootDirectory.Path(), "foo.txt"),
		})
	assert.Error(t, err, "only one source file allowed when destination is a file")
}

func TestCopyMultipleFilesWithNonExistingFile(t *testing.T) {
	rootDirectory := fs.NewDir(t, "root",
		fs.WithFile("foo.txt", "foo"),
		fs.WithFile("bar.txt", "bar"))
	defer rootDirectory.Remove()

	src := filepath.Join(rootDirectory.Path(), "source", "non-existing")
	err := copy(
		filepath.Join(rootDirectory.Path(), "destination"),
		[]string{
			src,
			filepath.Join(rootDirectory.Path(), "foo.txt"),
		})
	assert.Error(t, err, fmt.Sprintf("file %q does not exist", src))
}

func TestCopyMultipleFilesWithNonExistingGlob(t *testing.T) {
	rootDirectory := fs.NewDir(t, "root",
		fs.WithFile("foo.txt", "foo"),
		fs.WithFile("bar.txt", "bar"))
	defer rootDirectory.Remove()

	src := filepath.Join(rootDirectory.Path(), "source", "non-existing*")
	err := copy(
		filepath.Join(rootDirectory.Path(), "destination"),
		[]string{
			src,
			filepath.Join(rootDirectory.Path(), "foo.txt"),
		})
	assert.Error(t, err, fmt.Sprintf("file %q does not exist", src))
}

func TestCopyNonExistingFile(t *testing.T) {
	rootDirectory := fs.NewDir(t, "root")
	defer rootDirectory.Remove()

	err := copy(filepath.Join(rootDirectory.Path(), "bar.txt"),
		[]string{filepath.Join(rootDirectory.Path(), "non-existing")})
	assert.ErrorContains(t, err, "does not exist")
}

func TestCopyFileToNonExistingPathFile(t *testing.T) {
	rootDirectory := fs.NewDir(t, "root",
		fs.WithFile("foo.txt", "foo"))
	defer rootDirectory.Remove()

	err := copy(filepath.Join(rootDirectory.Path(), "bar", "bar.txt"),
		[]string{filepath.Join(rootDirectory.Path(), "foo.txt")})
	assert.NilError(t, err)

	expected := fs.Expected(t,
		fs.WithFile("foo.txt", "foo"),
		fs.WithDir("bar", fs.WithMode(0700),
			fs.WithFile("bar.txt", "foo")))
	assert.Assert(t, fs.Equal(rootDirectory.Path(), expected))
}

func TestCopyFileOverItself(t *testing.T) {
	rootDirectory := fs.NewDir(t, "root",
		fs.WithFile("foo.txt", "foo"))
	defer rootDirectory.Remove()

	f := filepath.Join(rootDirectory.Path(), "foo.txt")
	err := copy(f, []string{f})
	assert.NilError(t, err)

	expected := fs.Expected(t, fs.WithFile("foo.txt", "foo"))
	assert.Assert(t, fs.Equal(rootDirectory.Path(), expected))
}

func TestCopyTreeToNonExistingPath(t *testing.T) {
	rootDirectory := fs.NewDir(t, "root",
		fs.WithDir("source",
			fs.WithFile("foo.txt", "foo"),
			fs.WithDir("bar",
				fs.WithFile("bar.txt", "bar"))))
	defer rootDirectory.Remove()

	err := copy(filepath.Join(rootDirectory.Path(), "destination"),
		[]string{filepath.Join(rootDirectory.Path(), "source", "bar")})
	assert.NilError(t, err)

	expected := fs.Expected(t,
		fs.WithDir("source",
			fs.WithFile("foo.txt", "foo"),
			fs.WithDir("bar",
				fs.WithFile("bar.txt", "bar"))),
		fs.WithDir("destination",
			fs.WithDir("bar",
				fs.WithFile("bar.txt", "bar"))))
	assert.Assert(t, fs.Equal(rootDirectory.Path(), expected))
}

func testCopyDeepTreeToExistingPath(t *testing.T) {
	rootDirectory := fs.NewDir(t, "root",
		fs.WithDir("source",
			fs.WithFile("bar.txt", "bar"),
			fs.WithDir("subdir1",
				fs.WithFile("foo.txt", "foo"),
				fs.WithDir("subdir2",
					fs.WithFile("qix.txt", "qix"),
				),
			)),
		fs.WithDir("destination"))

	err := copy(filepath.Join(rootDirectory.Path(), "destination"),
		[]string{filepath.Join(rootDirectory.Path(), "source")})
	assert.NilError(t, err)

	expected := fs.Expected(t,
		fs.WithDir("source",
			fs.WithFile("bar.txt", "bar"),
			fs.WithDir("subdir1",
				fs.WithFile("foo.txt", "foo"),
				fs.WithDir("subdir2",
					fs.WithFile("qix.txt", "qix"),
				),
			)),
		fs.WithDir("destination"),
		fs.WithFile("bar.txt", "bar"),
		fs.WithDir("subdir1",
			fs.WithFile("foo.txt", "foo"),
			fs.WithDir("subdir2",
				fs.WithFile("qix.txt", "qix"),
			),
		))
	assert.Assert(t, fs.Equal(rootDirectory.Path(), expected))
}

func TestCopyTreeWithGlob(t *testing.T) {
	rootDirectory := fs.NewDir(t, "root",
		fs.WithDir("source",
			fs.WithFile("foo.txt", "foo"),
			fs.WithDir("bar",
				fs.WithFile("bar.txt", "bar"))))
	defer rootDirectory.Remove()

	err := copy(filepath.Join(rootDirectory.Path(), "destination"),
		[]string{filepath.Join(rootDirectory.Path(), "source", "*")})
	assert.NilError(t, err)

	expected := fs.Expected(t,
		fs.WithDir("source",
			fs.WithFile("foo.txt", "foo"),
			fs.WithDir("bar",
				fs.WithFile("bar.txt", "bar")),
		),
		fs.WithDir("destination",
			fs.WithFile("foo.txt", "foo"),
			fs.WithDir("bar",
				fs.WithFile("bar.txt", "bar"))))
	assert.Assert(t, fs.Equal(rootDirectory.Path(), expected))
}

func TestCopyTreeWithEmptyGlob(t *testing.T) {
	rootDirectory := fs.NewDir(t, "root",
		fs.WithDir("source",
			fs.WithFile("foo.txt", "foo")))
	defer rootDirectory.Remove()

	src := filepath.Join(rootDirectory.Path(), "source", "non-existing*")
	err := copy(filepath.Join(rootDirectory.Path(), "destination"), []string{src})
	assert.Error(t, err, fmt.Sprintf("file %q does not exist", src))

	expected := fs.Expected(t,
		fs.WithDir("source",
			fs.WithFile("foo.txt", "foo")))
	assert.Assert(t, fs.Equal(rootDirectory.Path(), expected))
}
