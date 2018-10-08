package build

import (
	"flag"

	"github.com/mat007/brique"
)

var (
	testRun = flag.String("test.run", "", "pattern to filter the tests")
)

// All does everything
func All(b *building.B) {
	Depends(b)
	Test(b)
}

// Test runs the tests
func Test(b *building.B) {
	b.Go("test", "-test.run", *testRun, "./...")
}

// Depends retrieves the dependencies
func Depends(b *building.B) {
	b.Dep("ensure")
}
