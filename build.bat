@echo off

rem set the project package name
set PACKAGE_NAME=github.com/mat007/b

rem configure current platform
set MSYS_NO_PATHCONV=1
set GOOS=windows

rem test if go is available
go version >nul 2>nul
if errorlevel 1 goto docker
rem build the build tool directly
go build ./cmd/b && b.exe %*
goto end

:docker
rem test if docker is available
docker version >nul 2>nul
if errorlevel 1 goto error
rem build the build tool in a container
docker run --rm -t -v%cd%:/go/src/%PACKAGE_NAME% -e GOOS=%GOOS% ^
    -w /go/src/%PACKAGE_NAME% golang:1.10.3-alpine3.7 go build ./cmd/b ^
    && b.exe %*
goto end

:error
echo Either go ^(http://golang.org^) or docker ^(http://www.docker.com^) needed to build.
exit /b 1

:end
