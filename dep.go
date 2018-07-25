package b

var depVersion = "v0.4.1"

func Dep(args ...string) Tool {
	t := MakeTool(
		"dep",
		"version",
		"https://github.com/golang/dep",
		`
FROM alpine:`+alpineVersion+`
RUN apk add --no-cache curl

RUN curl -o /usr/bin/dep -L https://github.com/golang/dep/releases/download/`+depVersion+`/dep-linux-amd64 && \
    chmod +x /usr/bin/dep
`)
	if len(args) > 0 {
		t.Run(args...)
	}
	return t
}
