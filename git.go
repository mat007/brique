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
FROM alpine:`+AlpineVersion+`
RUN apk add --no-cache git
`)
	})
	if len(args) > 0 {
		gitTool.Run(args...)
	}
	return gitTool
}

func GitShortCommit() string {
	buf := &bytes.Buffer{}
	Git().WithOutput(buf).WithSuccess().Run("rev-parse", "--short", "HEAD")
	return strings.TrimSuffix(buf.String(), "\n")
}

func GitCommit() string {
	buf := &bytes.Buffer{}
	Git().WithOutput(buf).WithSuccess().Run("rev-parse", "HEAD")
	return strings.TrimSuffix(buf.String(), "\n")
}

func GitTag() string {
	buf := &bytes.Buffer{}
	Git().WithOutput(buf).WithSuccess().Run("describe", "--always", "--dirty")
	return strings.TrimSuffix(buf.String(), "\n")
}

func GitDirty() bool {
	Git("update-index", "-q", "--refresh")
	return Git().WithSuccess().Run("diff-index", "--quiet", "HEAD", "--", ".")
}

func GitVersion() string {
	buf := &bytes.Buffer{}
	Git().WithOutput(buf).WithSuccess().Run("tag", "-l", "--points-at", "HEAD", `"v*"`)
	return strings.TrimSuffix(buf.String(), "\n")
}
