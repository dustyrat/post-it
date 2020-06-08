#!/bin/bash
ext() { 
    local array=$1 index=$2
    local i="${array}_$index"
    printf '%s' "${!i}"
}
declare "EXT_windows=.exe"

VERSION=$(git rev-parse HEAD)
TAG=$(git tag --points-at "$VERSION")
if [ -n "$TAG" ]; then 
  VERSION=$TAG
fi
echo VERSION: $VERSION 
for OS in "darwin" "linux" "windows"; do
    for ARCH in "386" "amd64"; do
        echo Building: release/${OS}_${ARCH}/post-it$(ext EXT $OS)
        CGO_ENABLED=0 GOOS=$OS GOARCH=$ARCH go build -ldflags "-X main.version=$VERSION" -o release/${OS}_${ARCH}/post-it$(ext EXT $OS) ./main.go
    done
done
