package b

import (
	"sync"
)

var (
	dockerTool Tool
	dockerOnce sync.Once
)

func Docker(args ...string) Tool {
	goOnce.Do(func() {
		// $$$$ MAT check what happens with empty instructions
		dockerTool = MakeTool("docker", "--version", "https://www.docker.com", "")
	})
	if len(args) > 0 {
		dockerTool.Run(args...)
	}
	return dockerTool
}
