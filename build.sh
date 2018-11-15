#!/bin/sh
set -e

# set the project package name
PACKAGE_NAME=github.com/mat007/brique

# configure the current platform
if [ "${OS}" = "Windows_NT" ]; then
    export MSYS_NO_PATHCONV=1
    GOOS=windows
    EXE=.exe
else
    GOOS=$(uname -s | tr '[:upper:]' '[:lower:]')
    EXE=
fi

# test if build tool is available
if [ ! -x "$(command -v ./b)" ]; then
    mkdir -p b_main
    function cleanup {
        rm -r b_main
    }
    trap cleanup EXIT
    echo "package main
import \"github.com/mat007/brique\"
func main() { building.Main() }" > b_main/main.go
    # test if go is available
    if [ -x "$(command -v go)" ]; then
        # build b
        go build -o b$EXE b_main/main.go
    else
        # test if docker is available
        if [ -x "$(command -v docker)" ]; then
            # build b in a container
            docker run --rm -t -v "$(pwd)":/go/src/$PACKAGE_NAME -e GOOS=$GOOS \
                -w /go/src/$PACKAGE_NAME golang:1.10.3-alpine3.7 go build -o b$EXE b_main/main.go
        else
            echo "Either Go (http://golang.org) or Docker (http://www.docker.com) is needed to build."
            exit 1
        fi
    fi
fi

# run the build tool forwarding the arguments
./b "$@"
