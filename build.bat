@echo off

rem set this to the build tool version
VERSION=v0.0.1

rem test if the build tool is available
if exist .\b.exe (
    .\b.exe %*
    exit /b %errorlevel%
)

rem test if docker is available
docker version >nul 2>nul
if errorlevel 1 (
    echo Either b ^(https://github.com/mat007/b^) or docker ^(http://www.docker.com^) needed to build.
    exit /b 1
)

rem download the build tool in a container
docker run --rm -t -v "%cd%":/pwd -w /pwd appropriate/curl curl --fail -o b.exe https://github.com/mat007/b/releases/download/$VERSION/b-windows.exe
if errorlevel 1 (
    exit /b %errorlevel%
)

.\b.exe %*
