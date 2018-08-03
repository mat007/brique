package building

func (b *B) Esc(args ...string) Tool {
	return b.MakeTool(
		"esc",
		"--help",
		"https://github.com/mjibson/esc",
		"FROM golang:"+GoVersion+"-alpine"+AlpineVersion+`
RUN apk add --no-cache git && \
	go get gopkg.in/mjibson/esc.v0 && \
	mv /go/bin/esc.v0 /go/bin/esc`,
		args...)
}
