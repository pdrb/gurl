#!/usr/bin/env bash

# Set bash strict mode
set -euo pipefail
IFS=$'\n\t'

# Binaries names
linux_amd64_file='gurl-linux-x86_64'
linux_arm64_file='gurl-linux-arm64'
mac_amd64_file='gurl-mac-x86_64'
mac_arm64_file='gurl-mac-arm64'

# Show go version
go version
echo

# Linux
echo "------- Linux Build -------"
echo "- Building for amd64..."
GOOS=linux GOARCH=amd64 go build -o "$linux_amd64_file" main.go
gzip -f "$linux_amd64_file"
echo "- Building for arm64..."
GOOS=linux GOARCH=arm64 go build -o "$linux_arm64_file" main.go
gzip -f "$linux_arm64_file"
echo

# Mac
echo "------- Mac Build -------"
echo "- Building for amd64..."
GOOS=darwin GOARCH=amd64 go build -o "$mac_amd64_file" main.go
gzip -f "$mac_amd64_file"
echo "- Building for arm64..."
GOOS=darwin GOARCH=arm64 go build -o "$mac_arm64_file" main.go
gzip -f "$mac_arm64_file"
echo

echo -e "Build successfully! :)"
