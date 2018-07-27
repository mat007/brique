package b

var GoMetaLinterVersion = "2.0.5"

func GoMetaLinter(args ...string) Tool {
	t := MakeTool(
		"gometalinter",
		"--version",
		"https://github.com/alecthomas/gometalinter",
		`FROM golang:`+GoVersion+`-alpine`+AlpineVersion+`
WORKDIR /go/src/github.com/alecthomas/gometalinter
RUN apk add --no-cache curl && \
    curl -L https://github.com/alecthomas/gometalinter/archive/v`+GoMetaLinterVersion+`.tar.gz | tar xz --strip-components=1 && \
	go build -v -o /usr/local/bin/gometalinter . && \
	gometalinter --install && \
	rm -rf /go/src/* /go/pkg/*`)
	if len(args) > 0 {
		t.Run(args...)
	}
	return t
}
