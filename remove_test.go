package building

import (
	"path/filepath"
	"testing"

	"gotest.tools/assert"
	"gotest.tools/fs"
)

func TestRemove(t *testing.T) {
	rootDirectory := fs.NewDir(t, "root",
		fs.WithDir("empty-dir"),
		fs.WithDir("full-dir",
			fs.WithFile("some-file", "")),
		fs.WithFile("file", ""),
		fs.WithFile("remaining-file", ""),
		fs.WithDir("remaining-dir"))
	defer rootDirectory.Remove()

	err := remove(
		filepath.Join(rootDirectory.Path(), "empty-dir"),
		filepath.Join(rootDirectory.Path(), "full-dir"),
		filepath.Join(rootDirectory.Path(), "non-existing"),
		filepath.Join(rootDirectory.Path(), "file"))
	assert.NilError(t, err)

	expected := fs.Expected(t,
		fs.WithFile("remaining-file", ""),
		fs.WithDir("remaining-dir"))
	assert.Assert(t, fs.Equal(rootDirectory.Path(), expected))
}

func TestRemoveWithGlob(t *testing.T) {
	rootDirectory := fs.NewDir(t, "root",
		fs.WithDir("empty-dir"),
		fs.WithDir("full-dir",
			fs.WithFile("some-file", "")),
		fs.WithFile("remaining-file", ""))
	defer rootDirectory.Remove()

	err := remove(
		filepath.Join(rootDirectory.Path(), "*-dir"),
		filepath.Join(rootDirectory.Path(), "full-dir"))
	assert.NilError(t, err)

	expected := fs.Expected(t, fs.WithFile("remaining-file", ""))
	assert.Assert(t, fs.Equal(rootDirectory.Path(), expected))
}

func TestRemoveWithNonExistingGlob(t *testing.T) {
	rootDirectory := fs.NewDir(t, "root",
		fs.WithFile("remaining-file", ""))
	defer rootDirectory.Remove()

	err := remove(filepath.Join(rootDirectory.Path(), "non-existing*"))
	assert.NilError(t, err)
	expected := fs.Expected(t, fs.WithFile("remaining-file", ""))
	assert.Assert(t, fs.Equal(rootDirectory.Path(), expected))
}
