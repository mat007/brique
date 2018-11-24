package building

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"gotest.tools/assert"
	"gotest.tools/fs"
)

func TestFilesetGlob(t *testing.T) {
	rootDirectory := fs.NewDir(t, "root",
		fs.WithFile("foo.txt", "foo"),
		fs.WithFile("bar.txt", "bar"),
		fs.WithDir("bar", fs.WithMode(0700),
			fs.WithFile("bla.txt", "bla"),
			fs.WithDir("sub", fs.WithMode(0700),
				fs.WithFile("mar.txt", "mar"))))
	defer rootDirectory.Remove()

	checkGlobFileset(t, rootDirectory, "", "",
		[]string{"."},
		[]string{".", "bar", "bar/bla.txt", "bar/sub", "bar/sub/mar.txt", "bar.txt", "foo.txt"}) // $$$$ MAT: test remove, copy, tar & zip with this fileset
	checkGlobFileset(t, rootDirectory, "bar/*", "",
		[]string{"bar/bla.txt", "bar/sub"},
		[]string{"bar/bla.txt", "bar/sub", "bar/sub/mar.txt"})
	checkGlobFileset(t, rootDirectory, filepath.Join(rootDirectory.Path(), "foo.txt"), "", []string{"foo.txt"})
	checkGlobFileset(t, rootDirectory, "*", "",
		[]string{"bar", "bar.txt", "foo.txt"},
		[]string{"bar", "bar/bla.txt", "bar/sub", "bar/sub/mar.txt", "bar.txt", "foo.txt"})
	checkGlobFileset(t, rootDirectory, "*,*", "",
		[]string{"bar", "bar.txt", "foo.txt", "bar", "bar.txt", "foo.txt"},
		[]string{"bar", "bar/bla.txt", "bar/sub", "bar/sub/mar.txt", "bar.txt", "foo.txt"})
	checkGlobFileset(t, rootDirectory, "ba*", "",
		[]string{"bar", "bar.txt"},
		[]string{"bar", "bar/bla.txt", "bar/sub", "bar/sub/mar.txt", "bar.txt"})
	checkGlobFileset(t, rootDirectory, "bar", "",
		[]string{"bar"},
		[]string{"bar", "bar/bla.txt", "bar/sub", "bar/sub/mar.txt"})
	checkGlobFileset(t, rootDirectory, "bar/*", "",
		[]string{"bar/bla.txt", "bar/sub"},
		[]string{"bar/bla.txt", "bar/sub", "bar/sub/mar.txt"})
	checkGlobFileset(t, rootDirectory, "*/sub", "",
		[]string{"bar/sub"},
		[]string{"bar/sub", "bar/sub/mar.txt"})
	checkGlobFileset(t, rootDirectory, "*/sub/*", "", []string{"bar/sub/mar.txt"})

	checkGlobFileset(t, nil, "", "", []string{"."}, nil) // $$$$ MAT: test remove, copy, tar & zip with this fileset
	checkGlobFileset(t, nil, filepath.Join(rootDirectory.Path(), "foo.txt"), "",
		[]string{filepath.ToSlash(filepath.Join(rootDirectory.Path(), "foo.txt"))},
		[]string{"foo.txt"})
	checkGlobFileset(t, nil, filepath.Join(rootDirectory.Path(), "*"), "",
		[]string{
			filepath.ToSlash(filepath.Join(rootDirectory.Path(), "bar")),
			filepath.ToSlash(filepath.Join(rootDirectory.Path(), "bar.txt")),
			filepath.ToSlash(filepath.Join(rootDirectory.Path(), "foo.txt"))},
		[]string{"bar", "bar/bla.txt", "bar/sub", "bar/sub/mar.txt", "bar.txt", "foo.txt"})

	checkGlobFileset(t, rootDirectory, "*", "bar",
		[]string{"bar", "bar.txt", "foo.txt"},
		[]string{"bar.txt", "foo.txt"})
	checkGlobFileset(t, rootDirectory, "*", "bar/*",
		[]string{"bar", "bar.txt", "foo.txt"})
	checkGlobFileset(t, rootDirectory, "*", "bar*",
		[]string{"bar", "bar.txt", "foo.txt"},
		[]string{"foo.txt"})
	checkGlobFileset(t, rootDirectory, "*", "*/sub",
		[]string{"bar", "bar.txt", "foo.txt"},
		[]string{"bar", "bar/bla.txt", "bar.txt", "foo.txt"})
	checkGlobFileset(t, rootDirectory, "*", "*/sub/*",
		[]string{"bar", "bar.txt", "foo.txt"},
		[]string{"bar", "bar/bla.txt", "bar/sub", "bar.txt", "foo.txt"})
}

func checkGlobFileset(t *testing.T, dir *fs.Dir, includes, excludes string, expected ...[]string) {
	t.Helper()
	fs := fileset{
		includes: strings.Split(includes, ","),
		excludes: []string{excludes},
	}
	if dir != nil {
		fs.dir = dir.Path()
	}
	actual, err := fs.glob(false)
	assert.NilError(t, err)
	if dir != nil {
		assert.Equal(t, actual.dir, dir.Path())
	} else {
		assert.Equal(t, actual.dir, "")
	}
	assert.DeepEqual(t, actual.includes, expected[0])
	expectedWalk := expected[0]
	if len(expected) > 1 {
		expectedWalk = expected[1]
	}
	if expectedWalk == nil {
		return
	}
	var actualWalk []string
	err = actual.walk(func(path, rel string, info os.FileInfo, err error) error {
		actualWalk = append(actualWalk, rel)
		return err
	})
	assert.NilError(t, err)
	assert.DeepEqual(t, actualWalk, expectedWalk)
}
