package b

func Jq(args ...string) Tool {
	t := MakeTool(
		"jq",
		"--help",
		"https://stedolan.github.io/jq",
		`
FROM alpine:`+AlpineVersion+`
RUN apk add --no-cache jq
`)
	if len(args) > 0 {
		t.Run(args...)
	}
	return t
}
