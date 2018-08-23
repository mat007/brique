#!/bin/sh -e

# set the project package name
PACKAGE_NAME=github.com/mat007/b
# in a user project this would point to the vendored folder
BUILD_TOOL=./cmd/b

# configure current platform
if [ "${OS}" = "Windows_NT" ]; then
    export MSYS_NO_PATHCONV=1
    GOOS=windows
else
    GOOS=$(uname -s | tr '[:upper:]' '[:lower:]')
fi

# test if go is available
if [ -x "$(command -v go)" ]; then
    # build the build tool directly
    go build $BUILD_TOOL
else
    # test if docker is available
    if [ -x "$(command -v docker)" ]; then
        # build the build tool in a container
        docker run --rm -t -v"$(pwd)":/go/src/$PACKAGE_NAME -e GOOS=$GOOS \
            -w /go/src/$PACKAGE_NAME golang:1.10.3-alpine3.7 go build $BUILD_TOOL
    else
        echo "Either go (http://golang.org) or docker (http://www.docker.com) needed to build."
        exit 1
    fi
fi

# run the build tool forwarding the arguments
./b $*
