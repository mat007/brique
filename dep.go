package building

var DepVersion = "v0.4.1"

func (b *B) Dep(args ...string) Tool {
	return b.MakeTool(
		"dep",
		"version",
		"https://github.com/golang/dep",
		"FROM golang:"+GoVersion+"-alpine"+AlpineVersion+`
RUN apk add --no-cache git curl && \
    curl -o /usr/bin/dep -L https://github.com/golang/dep/releases/download/`+DepVersion+`/dep-linux-amd64 && \
    chmod +x /usr/bin/dep`,
		args...)
}
