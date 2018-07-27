package b

import (
	"flag"
)

var (
	cross    = flag.Bool("cross", false, "build for all platforms (linux, darwin, windows)")
	parallel = flag.Bool("parallel", false, "build in parallel")

	GoVersion = "1.10.3"
)

func Go(args ...string) Tool {
	return MakeTool(
		"go",
		"version",
		"http://golang.org",
		"FROM golang:"+GoVersion+"-alpine"+AlpineVersion,
		args...)
}
