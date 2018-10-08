# Brique

A quest for the perfect build tool.

## Why Brique?

Many Go projects use `make` to drive the build, however it has several drawbacks including its severe lack of portability on the Windows platform.

Even if it's possible to find an implementation of `make` for Windows, it's not as readily available as for Linux and Darwin.
The biggest problem however is that a `Makefile` usually relies on the underlying shell for most features and something as simple as `tar` or even `rm` is difficult to provide or write properly on Windows.

For any good size multi-platform project this often means ending up with build code difficult to maintain sometimes using multiple scripting languages (`Makefile`, `shell`, `bash`, `PowerShell`, `batch`, `python`, …).

Brique aims at unifying all this by making it possible to use Go as the build code.

## What is Brique?

Brique is a build tool for Go projects made of three parts:
* a Go binary called `b` (or `b.exe` on Windows) which creates and invokes on the fly a custom binary specifically tailored to build each project
* a Go package of bricks (called `building`) to help with implementing the build features needed by each project
* a Shell and a Batch script to bootstrap the build and allow to vendor the whole tool chain in each project

The goal of Brique is to create a build environment which is:
* with a low barrier of entry (build written in Go, no tool chain to install)
* fast (especially when nothing needs to be done)
* extensible (linters, coverage, release, signing, …)
* multi-platform (single code, cross-compilation, parallel builds)
* build server friendly (reproducible, cleans after itself, no requirements other than `Docker`)

Using containers provides a nice way to fulfill most of these requirements, however they tend to be slow:  either project files must be baked into an image and artifacts copied back from a container, or volumes must be mounted to share project files with a container and they are quite slow (on Windows mainly and Darwin also).

Brique works around this by providing build fallbacks.
For instance if a tool (e.g. `go`) is available it will be used directly, if not but `Docker` is available, Brique will spin a container to use the tool.

This flexibility allows for casual developpers on a project (or product managers or build servers) to build a project right away without having to figure out which tool chain is needed, while core developpers building more frequently will probably want to install most of the required tools to minimize the build times.

Brique provides a library of build tools wrappers and is very easy to extend with more custom wrappers as needed.

## How to get started with Brique?

Of course Brique can be used on itself!

For a quick glance at a project build code using Brique look at [cmd/build/build.go](cmd/build/build.go).

To find out how to build simply invoke `./build.sh` or `build.bat` at the root of the project with the `-help` flag:
```
$ ./build.sh -help
Usage: build [OPTIONS] [TARGETS]

Options:
  -containers
        always build in containers
  -cross
        build for all platforms (linux, darwin, windows)
  -parallel
        build in parallel
  -q    quiet output
  -test.run string
        pattern to filter the tests
  -v    verbose output

Targets:
  all   does everything
  depends
        retrieves the dependencies
  test  runs the tests
```

This output gets generated on the fly, a bit like what `go test` does with parsing files ending in `_test.go` to extract test functions matching the test designated signature, except here Brique looks for an exported `func Xyz(b *building.B)`.

The comments above target functions are used as descriptions for the targets listed in the help.

The first target in the Go build file becomes the default one, meaning here calling `./build.sh` or `build.bat` with no argument will invoke the `all` target, e.g.:
```
$ ./build.sh
build started
> all
running [dep ensure]
running [go test -test.run  ./...]
ok      github.com/mat007/brique        (cached)
?       github.com/mat007/brique/cmd/build      [no test files]
< all (took 19.5013775s)
build finished (took 19.503879s)
```

## How to use Brique?

Brique releases no binary because it's designed to be entirely vendored.

To add Brique to a project follow these steps:
* Add both [build.sh](build.sh) and [build.bat](build.bat) to the project root folder
* In both `build.sh` and `build.bat` change the `PACKAGE_NAME` variable to the project package name
* Create a file `cmd/build/build.go` with the content:
```go
package build

import (
	"github.com/mat007/brique"
)

// All does everything
func All(b *building.B) {
	fmt.Println("OK!)
}
```
* Use your vendoring tool to bring Brique in
* Test by invoking `build.sh` or `build.bat` from the project root folder

## Where to go next with Brique?

The project is still in alpha phase and has no stable API yet, meaning it may be a bit early to use in production.
However the only way to gather feedback and shape it properly now would be to start using it in real world projects.

It also needs proper end-to-end tests, continuous integration, more documentation and also feedback, feedback and feedback!
