package b

import (
	"bytes"
	"strings"
	"sync"
)

var (
	gitTool Tool
	gitOnce sync.Once
)

func Git(args ...string) Tool {
	gitOnce.Do(func() {
		gitTool = MakeTool(
			"git",
			"--version",
			"https://git-scm.com/",
			`
FROM alpine:`+alpineVersion+`
RUN apk add --no-cache git
`)
	})
	if len(args) > 0 {
		gitTool.Run(args...)
	}
	return gitTool
}

func GitCommit() string {
	buf := &bytes.Buffer{}
	Git().WithOutput(buf).WithSuccess().Run("rev-parse", "--short", "HEAD")
	return strings.TrimSuffix(buf.String(), "\n")
}

func GitTag() string {
	buf := &bytes.Buffer{}
	Git().WithOutput(buf).WithSuccess().Run("describe", "--always", "--dirty")
	return strings.TrimSuffix(buf.String(), "\n")
}
