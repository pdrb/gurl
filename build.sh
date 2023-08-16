#!/usr/bin/env bash

# Set bash strict mode
set -euo pipefail
IFS=$'\n\t'

# Binaries names
linux_amd64_file='gurl-linux-x86_64'
linux_arm64_file='gurl-linux-arm64'
mac_amd64_file='gurl-mac-x86_64'
mac_arm64_file='gurl-mac-arm64'
windows_amd64_file='gurl-windows-x86_64.exe'
windows_arm64_file='gurl-windows-arm64.exe'

# Show go version
go version
echo

# Linux
echo "------- Linux Build -------"
echo "- Building for amd64..."
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o "$linux_amd64_file"
gzip -f "$linux_amd64_file"
echo "- Building for arm64..."
CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -ldflags="-s -w" -o "$linux_arm64_file"
gzip -f "$linux_arm64_file"
echo

# Mac
echo "------- Mac Build -------"
echo "- Building for amd64..."
CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -ldflags="-s -w" -o "$mac_amd64_file"
gzip -f "$mac_amd64_file"
echo "- Building for arm64..."
CGO_ENABLED=0 GOOS=darwin GOARCH=arm64 go build -ldflags="-s -w" -o "$mac_arm64_file"
gzip -f "$mac_arm64_file"
echo

# Windows
echo "------- Windows Build -------"
echo "- Building for amd64..."
CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -ldflags="-s -w" -o "$windows_amd64_file"
gzip -f "$windows_amd64_file"
echo "- Building for arm64..."
CGO_ENABLED=0 GOOS=windows GOARCH=arm64 go build -ldflags="-s -w" -o "$windows_arm64_file"
gzip -f "$windows_arm64_file"
echo

echo -e "Build successfully! :)"
