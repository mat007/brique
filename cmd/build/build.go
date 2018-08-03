package build

import (
	"flag"
	"fmt"
	"os"

	"github.com/docker/cli/cli/command"
	"github.com/docker/cli/cli/flags"
	"github.com/mat007/b"
	"github.com/mat007/b/cmd/build/internal"
)

var (
	testF   = flag.String("test.run", "", "pattern to filter the tests")
	binName = flag.String("bin.name", "b", "name of the executable")
)

// TargetAll does everything
func TargetAll(b *building.B) {
	TargetTest(b)
	TargetBin(b)
}

// TargetBin builds the binaries
func TargetBin(b *building.B) {
	b.WithOS(func(goos string) {
		b.Go().WithEnv("GOOS="+goos, "CGO_ENABLED=0").
			Run("build", "-ldflags=-s -w", "-o", *binName+"-"+goos+b.ExecExt(goos), "./cmd/b")
	})
}

// TargetTest runs the tests
func TargetTest(b *building.B) {
	b.Go("test", "-test.run", *testF, "./...")
}

// TargetDepends retrieves the dependencies
func TargetDepends(b *building.B) {
	b.Dep("ensure")
}

// TargetClean cleans the build artifacts
func TargetClean(b *building.B) {
	b.Remove(*binName + "-*")
}

// TargetHello demos the build tool
func TargetHello(b *building.B) {
	fmt.Println("Hello !")
	print.Hello()
	dockerCli := command.NewDockerCli(os.Stdin, os.Stdout, os.Stderr, false)
	opts := flags.NewClientOptions()
	dockerCli.Initialize(opts)
	fmt.Println("docker:", dockerCli.DefaultVersion())
}
