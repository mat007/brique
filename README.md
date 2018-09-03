# Brique

A quest for the perfect build tool.

## Why Brique?

Many Go projects use `make` to drive the build, however it has several drawbacks including its severe lack of portability on the Windows platform.

Even if it's possible to find an implementation of `make` for Windows, it's not as readily available as for Linux and Darwin.
The biggest problem however is that a `Makefile` usually relies on the underlying shell for most features and something as simple as `tar` or even `rm` is difficult to provide or write properly on Windows.

For any good size multi-platform project this often means ending up with build code difficult to maintain sometimes using multiple scripting languages (`Makefile`, `shell`, `bash`, `PowerShell`, `batch`, `python`, …).

Brique aims at unifying all this by making it possible to use Go as the build code.

## What is Brique?

Brique is a build tool for Go projects made of two parts:
* a Go binary called `b` (or `b.exe` on Windows) which creates and invokes on the fly a custom binary specifically tailored to build each project
* a Go library of bricks as a `building` package to help with implementing the build features needed by each project

The goal of Brique is to create a build environment which is:
* multi-platform (single code, cross-compilation, parallel builds)
* with a low barrier of entry (build written in Go, no tool chain to install)
* build server friendly (reproducible, cleans after itself, no requirements other than `Docker`)
* extensible (linters, coverage, release, signing, …)
* fast (especially when nothing needs to be done)

Using containers provides a nice way to fulfill most of these requirements, however they tend to be slow:  either project files must be baked into an image and artifacts copied back from a container, or volumes must be mounted to share project files with a container and they are quite slow (on Windows mainly and Darwin also).

Brique works around this by providing build fallbacks.
For instance if a tool (e.g. `go`) is available it will be used directly, if not but `Docker` is available, Brique will spin a container to use the tool.

This flexibility allows for casual developpers on a project (or product managers or build servers) to build a project right away without having to figure out which tool chain is needed, while core developpers building more frequently will probably want to install most of the required tools to minimize the build times.

Brique provides a library of build tools wrappers and is very easy to extend with more custom wrappers as needed.

## How to get started with Brique?

Of course Brique builds itself!

For a quick glance at a project build code using Brique look at [cmd/build/build.go](cmd/build/build.go).

To find out how to build simply invoke `b` at the root of the project with the `-help` flag:
```
$ b -help
Usage: build [OPTIONS] [TARGETS]

Options:
  -bin.name string
        name of the executable (default "b")
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
  -version string
        version of the executable (default "774d33e7c721570a2618b5fe2a45733a68e03e44")

Targets:
  all   does everything
  bin   builds the binaries
  clean cleans the build artifacts
  depends
        retrieves the dependencies
  test  runs the tests
  ```

This output gets generated on the fly, a bit like what `go test` does with parsing files ending in `_test.go` to extract test functions matching the test designated signature, except here `b` looks for `func TargetXyz(b *building.B)`.

The comments above target functions are used as descriptions for the targets listed in the help.

The first target in the Go build file becomes the default one, meaning here calling `b` with no argument will invoke the `all` target, e.g. on Windows:
```
$ b
build started
> all
running [go test -test.run  ./...]
ok      github.com/mat007/b     (cached)
ok      github.com/mat007/b/cmd/b       (cached)
?       github.com/mat007/b/cmd/build   [no test files]
building for windows
running [go build -ldflags=-s -w -o b-windows.exe ./cmd/b]
< all (took 3.5919863s)
build finished (took 3.5929954s)
```

## Where to go next with Brique?

The project is still in alpha phase and has no stable API yet, meaning it may be a bit early to use in production.
However the only way to gather feedback and shape it properly now would be to start using it in real world projects.

It also needs proper end-to-end tests, continuous integration, more documentation and also feedback, feedback and feedback!
