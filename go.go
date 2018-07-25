package b

import (
	"flag"
	"sync"
)

var (
	cross    = flag.Bool("cross", false, "build for all platforms (linux, darwin, windows)")
	parallel = flag.Bool("parallel", false, "build in parallel")

	GoVersion = "1.10.3"
)

var (
	goTool Tool
	goOnce sync.Once
)

func Go(args ...string) Tool {
	goOnce.Do(func() {
		goTool = MakeTool(
			"go",
			"version",
			"http://golang.org",
			"FROM golang:"+GoVersion+"-alpine"+AlpineVersion)
	})
	if len(args) > 0 {
		goTool.Run(args...)
	}
	return goTool
}
