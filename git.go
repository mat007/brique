package building

import (
	"bytes"
	"strings"
)

func (b *B) Git(args ...string) Tool {
	return b.MakeTool(
		"git",
		"--version",
		"https://git-scm.com",
		`
FROM alpine:`+AlpineVersion+`
RUN apk add --no-cache git`,
		args...)
}

func (b *B) GitShortCommit() string {
	buf := &bytes.Buffer{}
	b.Git().WithOutput(buf).WithSuccess().Run("rev-parse", "--short", "HEAD")
	return strings.TrimSpace(buf.String())
}

func (b *B) GitCommit() string {
	buf := &bytes.Buffer{}
	b.Git().WithOutput(buf).WithSuccess().Run("rev-parse", "HEAD")
	return strings.TrimSpace(buf.String())
}

func (b *B) GitTag() string {
	buf := &bytes.Buffer{}
	b.Git().WithOutput(buf).WithSuccess().Run("describe", "--always", "--dirty")
	return strings.TrimSpace(buf.String())
}

func (b *B) GitDirty() bool {
	b.Git("update-index", "-q", "--refresh")
	return b.Git().WithSuccess().Run("diff-index", "--quiet", "HEAD", "--", ".") == 0
}

func (b *B) GitVersion() string {
	buf := &bytes.Buffer{}
	b.Git().WithOutput(buf).WithSuccess().Run("tag", "-l", "--points-at", "HEAD", `"v*"`)
	return strings.TrimSpace(buf.String())
}
