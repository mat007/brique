package b

var goMetaLinterVersion = "2.0.5"

func GoMetaLinter(args ...string) Tool {
	t := MakeTool(
		"gometalinter",
		"--version",
		"https://github.com/alecthomas/gometalinter",
		`
FROM golang:`+goVersion+`-alpine`+alpineVersion+`
RUN apk add --no-cache curl

WORKDIR /go/src/github.com/alecthomas/gometalinter
RUN curl -L https://github.com/alecthomas/gometalinter/archive/v`+goMetaLinterVersion+`.tar.gz | tar xz --strip-components=1 && \
	go build -v -o /usr/local/bin/gometalinter . && \
	gometalinter --install && \
	rm -rf /go/src/* /go/pkg/*
`)
	if len(args) > 0 {
		t.Run(args...)
	}
	return t
}
