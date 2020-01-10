#!/usr/bin/env bash

# Remove old binaries (if any)
rm -rf dist

# Build Linux binaries:
env GOOS=linux GOARCH=386 go build -ldflags="-s -w" -o "dist/lietuviskatv_linux_i386" *.go             # Linux i386
env GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o "dist/lietuviskatv_linux_x86_64" *.go         # Linux 64bit
env GOOS=linux GOARCH=arm GOARM=5 go build -ldflags="-s -w" -o "dist/lietuviskatv_linux_arm" *.go      # Linux armv5/armel/arm (it also works on armv6)
env GOOS=linux GOARCH=arm GOARM=7 go build -ldflags="-s -w" -o "dist/lietuviskatv_linux_armhf" *.go    # Linux armv7/armhf
env GOOS=linux GOARCH=arm64 go build -ldflags="-s -w" -o "dist/lietuviskatv_linux_aarch64" *.go        # Linux armv8/aarch64

# Build FreeBSD binary:
env GOOS=freebsd GOARCH=amd64 go build -ldflags="-s -w" -o "dist/lietuviskatv_freebsd_x86_64" *.go     # FreeBSD 64bit

# Build MacOS binary:
env GOOS=darwin GOARCH=amd64 go build -ldflags="-s -w" -o "dist/lietuviskatv_darwin_x86_64" *.go       # Darwin 64bit

# Build Windows binaries:
env GOOS=windows GOARCH=386 go build -ldflags="-s -w" -o "dist/lietuviskatv_windows_i386.exe" *.go     # Windows 32bit
env GOOS=windows GOARCH=amd64 go build -ldflags="-s -w" -o "dist/lietuviskatv_windows_x86_64.exe" *.go # Windows 64bit

# Compress binaries (risk of binary not working at all on some platforms):
# upx --best dist/*
