#!/bin/sh -e

# set this to the build tool version
VERSION=v0.0.1

# configure current platform
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
    # test if curl is available ?
    if [ -x "$(command -v curl)" ]; then
        curl --fail -o b$EXE https://github.com/mat007/brique/releases/download/$VERSION/b-$GOOS$EXE
    else
        # test if docker is available
        if [ -x "$(command -v docker)" ]; then
            # download the build tool in a container
            docker run --rm -t -v "$(pwd)":/pwd \
                -w /pwd appropriate/curl curl --fail -o b$EXE https://github.com/mat007/brique/releases/download/$VERSION/b-$GOOS$EXE
        else
            echo "Either b (https://github.com/mat007/brique), curl or docker (http://www.docker.com) is needed to build."
            exit 1
        fi
    fi
fi
# run the build tool forwarding the arguments
./b $*
