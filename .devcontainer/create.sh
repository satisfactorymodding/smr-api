#!/usr/bin/env bash
set -ex

# Install protoc (used by go:generate) and add it to path
# Not done in Dockerfile because it seems to lock out writing to some files go needs to write to?
version=25.4
PB_REL="https://github.com/protocolbuffers/protobuf/releases"
curl -LO $PB_REL/download/v$version/protoc-$version-linux-x86_64.zip
unzip protoc-$version-linux-x86_64.zip -d $HOME/.local/
export PATH="$PATH:$HOME/.local/bin"
go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.34.2
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.3.0
export PATH="$PATH:$(go env GOPATH)/bin"

# Install libwebp (otherwise go:generate tries and fails to build it itself or something)
# https://developers.google.com/speed/webp/docs/precompiled
libwebp_version=1.4.0
libwebp_file="libwebp-$libwebp_version-linux-x86-64.tar.gz"
curl -LO "https://storage.googleapis.com/downloads.webmproject.org/releases/webp/$libwebp_file"
mkdir -p $HOME/libwebp
tar xzvf $libwebp_file -C $HOME/libwebp
export PATH="$PATH:$HOME/libwebp/bin"

# Add as a safe git directory
git config --global --add safe.directory "/workspaces/smr_api"
