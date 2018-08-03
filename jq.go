package building

func (b *B) Jq(args ...string) Tool {
	return b.MakeTool(
		"jq",
		"--help",
		"https://stedolan.github.io/jq",
		`
FROM alpine:`+AlpineVersion+`
RUN apk add --no-cache jq`,
		args...)
}
