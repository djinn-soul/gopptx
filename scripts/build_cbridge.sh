#!/bin/bash
# scripts/build_cbridge.sh
set -e

outDir="bindings/c/build"
mkdir -p "$outDir"

libName="libgopptx.so"
if [[ "$OSTYPE" == "darwin"* ]]; then
    libName="libgopptx.dylib"
fi
pythonLibPath="python/gopptx/$libName"

echo "Building gopptx shared library..."
if [[ "${GOPPTX_RELEASE_BUILD:-}" == "1" || "${GOPPTX_RELEASE_BUILD:-}" == "true" ]]; then
    go build -trimpath -buildvcs=false -ldflags="-s -w" -o "$outDir/$libName" -buildmode=c-shared bindings/c/bridge.go
else
    go build -o "$outDir/$libName" -buildmode=c-shared bindings/c/bridge.go
fi

echo "Build successful!"
echo "Library: $outDir/$libName"
cp "$outDir/$libName" "$pythonLibPath"
echo "Python package library synced: $pythonLibPath"
echo "Header: $outDir/libgopptx.h"
