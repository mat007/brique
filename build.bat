@echo off

rem set the project package name
set PACKAGE_NAME=github.com/mat007/brique

rem test if the build tool is available
if exist .\b.exe (
    goto :run
)

md b_main >nul 2>nul
echo package main>b_main\main.go
echo import "github.com/mat007/brique">>b_main\main.go
echo func main() { building.Main() }>>b_main\main.go

rem test if go is available
go version >nul 2>nul
if errorlevel 1 (
    rem test if docker is available
    docker version >nul 2>nul
    if errorlevel 1 (
        echo Either Go ^(http://golang.org^) or Docker ^(http://www.docker.com^) is needed to build.
        goto :error
    ) else (
        rem build b in a container
        docker run --rm -t -v "%cd%":/go/src/%PACKAGE_NAME% -e GOOS=windows^
            -w /go/src/%PACKAGE_NAME% golang:1.10.3-alpine3.7 go build -o b.exe b_main/main.go
        if errorlevel 1 (
            goto :error
        )
    )
) else (
    rem build b
    go build -o b.exe b_main/main.go
    if errorlevel 1 (
        goto :error
    )
)
goto :run

:error
rd /S /Q b_main >nul 2>nul
exit /b 1

:run
rd /S /Q b_main >nul 2>nul
.\b.exe %*
