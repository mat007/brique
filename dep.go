package b

var DepVersion = "v0.4.1"

func Dep(args ...string) Tool {
	return MakeTool(
		"dep",
		"version",
		"https://github.com/golang/dep",
		`
FROM alpine:`+AlpineVersion+`
RUN apk add --no-cache curl && \
    curl -o /usr/bin/dep -L https://github.com/golang/dep/releases/download/`+DepVersion+`/dep-linux-amd64 && \
    chmod +x /usr/bin/dep`,
		args...)
}
