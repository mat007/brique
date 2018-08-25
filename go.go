package building

var GoVersion = "1.10.3"

// $$$$ MAT go verbose with -v ?
func (b *B) Go(args ...string) Tool {
	return b.MakeTool(
		"go",
		"version",
		"http://golang.org",
		"FROM golang:"+GoVersion+"-alpine"+AlpineVersion,
		args...)
}
