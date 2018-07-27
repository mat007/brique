package b

func Jq(args ...string) Tool {
	return MakeTool(
		"jq",
		"--help",
		"https://stedolan.github.io/jq",
		`
FROM alpine:`+AlpineVersion+`
RUN apk add --no-cache jq`,
		args...)
}
