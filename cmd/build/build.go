package build

import (
	"flag"

	"github.com/mat007/brique"
)

var (
	b = building.Builder()

	output  = flag.String("o", "b", "name of the produced binary")
	testRun = flag.String("test.run", "", "pattern to filter the tests")
	version = flag.String("version", b.GitCommit(), "version of the binary")
)

// All does everything
func All(b *building.B) {
	Test(b)
	Bin(b)
}

// Bin builds the binaries
func Bin(b *building.B) {
	b.WithOS(func(goos string) {
		b.Go().WithEnv("GOOS="+goos, "CGO_ENABLED=0").
			Run("build", "-ldflags=-s -w", "-o", *output+"-"+goos+b.Exe(goos), "./cmd/b")
	})
}

// Test runs the tests
func Test(b *building.B) {
	b.Go("test", "-test.run", *testRun, "./...")
}

// Depends retrieves the dependencies
func Depends(b *building.B) {
	b.Dep("ensure")
}

// Clean cleans the build artifacts
func Clean(b *building.B) {
	b.Remove(*output + "-*")
}
