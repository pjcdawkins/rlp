#!/usr/bin/env bash

set -e

filename=rlp

export GOOS=linux
export GOARCH=amd64

echo "Building for $GOOS/$GOARCH and outputting to file: $filename" >&2

go build -ldflags="-s -w" -o "$filename"
