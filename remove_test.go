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

	err := Remove{}.WithFiles(
		filepath.Join(rootDirectory.Path(), "empty-dir"),
		filepath.Join(rootDirectory.Path(), "full-dir"),
		filepath.Join(rootDirectory.Path(), "non-existing"),
		filepath.Join(rootDirectory.Path(), "file")).run()
	assert.NilError(t, err)

	expected := fs.Expected(t,
		fs.WithFile("remaining-file", ""),
		fs.WithDir("remaining-dir"))
	assert.Assert(t, fs.Equal(rootDirectory.Path(), expected))
}

func TestRemoveWithFiles(t *testing.T) {
	rootDirectory := fs.NewDir(t, "root",
		fs.WithDir("empty-dir"),
		fs.WithDir("full-dir",
			fs.WithFile("some-file", "")),
		fs.WithFile("remaining-file", ""))
	defer rootDirectory.Remove()

	err := Remove{}.WithFiles(
		filepath.Join(rootDirectory.Path(), "*-dir"),
		filepath.Join(rootDirectory.Path(), "full-dir")).run()
	assert.NilError(t, err)

	expected := fs.Expected(t, fs.WithFile("remaining-file", ""))
	assert.Assert(t, fs.Equal(rootDirectory.Path(), expected))
}

func TestRemoveWithNonExistingFiles(t *testing.T) {
	rootDirectory := fs.NewDir(t, "root",
		fs.WithFile("remaining-file", ""))
	defer rootDirectory.Remove()

	err := Remove{}.WithFiles(
		filepath.Join(rootDirectory.Path(), "non-existing*")).run()
	assert.NilError(t, err)
	expected := fs.Expected(t, fs.WithFile("remaining-file", ""))
	assert.Assert(t, fs.Equal(rootDirectory.Path(), expected))
}

func TestRemoveWithFileset(t *testing.T) {
	rootDirectory := fs.NewDir(t, "root",
		fs.WithDir("empty-dir"),
		fs.WithDir("full-dir",
			fs.WithFile("some-file", "")),
		fs.WithFile("remaining-file", ""))
	defer rootDirectory.Remove()

	err := Remove{}.WithFileset(rootDirectory.Path(), "*-dir,full-dir", "").run()
	assert.NilError(t, err)

	expected := fs.Expected(t, fs.WithFile("remaining-file", ""))
	assert.Assert(t, fs.Equal(rootDirectory.Path(), expected))
}

func TestRemoveWithNonExistingFileset(t *testing.T) {
	rootDirectory := fs.NewDir(t, "root",
		fs.WithFile("remaining-file", ""))
	defer rootDirectory.Remove()

	err := Remove{}.WithFileset(rootDirectory.Path(), "non-existing*", "").run()
	assert.NilError(t, err)
	expected := fs.Expected(t, fs.WithFile("remaining-file", ""))
	assert.Assert(t, fs.Equal(rootDirectory.Path(), expected))
}

func TestRemoveWithExcludes(t *testing.T) {
	rootDirectory := fs.NewDir(t, "root",
		fs.WithDir("empty-dir"),
		fs.WithDir("full-dir",
			fs.WithFile("some-file", "")),
		fs.WithFile("remaining-file", ""))
	defer rootDirectory.Remove()

	err := Remove{}.WithFileset(rootDirectory.Path(), "*", "*-dir").run()
	assert.NilError(t, err)

	expected := fs.Expected(t,
		fs.WithDir("empty-dir"),
		fs.WithDir("full-dir",
			fs.WithFile("some-file", "")))
	assert.Assert(t, fs.Equal(rootDirectory.Path(), expected))
}
