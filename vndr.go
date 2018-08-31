package building

var VndrCommit = "1fc68ee0c852556a9ed53cbde16247033f104111"

// $$$$ MAT: with-git-proxy
func (b *B) Vndr(args ...string) Tool {
	return b.MakeTool(
		"vndr",
		"--help",
		"https://github.com/LK4D4/vndr", `
FROM golang:`+GoVersion+"-alpine"+AlpineVersion+`
RUN apk add --no-cache git \
 && mkdir -p $GOPATH/src/github.com/LK4D4 \
 && cd $GOPATH/src/github.com/LK4D4 \
 && git clone -q https://github.com/LK4D4/vndr \
 && cd $GOPATH/src/github.com/LK4D4/vndr \
 && git checkout -q `+VndrCommit+` \
 && go build -o /go/bin/vndr github.com/LK4D4/vndr`,
		args...)
}
