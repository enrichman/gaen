#!/bin/bash

VERSION=$(git describe --tags)
ARCH=amd64

for OS in linux darwin windows;
do
    GOARCH=${ARCH} GOOS=${OS} go build -ldflags "-X main.version=${VERSION}" -o build/gaen-${VERSION}-${OS}-${ARCH} .
done
