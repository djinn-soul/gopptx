#!/bin/bash
# scripts/build_cbridge.sh
set -e

outDir="bindings/c/build"
mkdir -p "$outDir"

libName="libgopptx.so"
if [[ "$OSTYPE" == "darwin"* ]]; then
    libName="libgopptx.dylib"
fi

echo "Building gopptx shared library..."
go build -o "$outDir/$libName" -buildmode=c-shared bindings/c/bridge.go

echo "Build successful!"
echo "Library: $outDir/$libName"
echo "Header: $outDir/libgopptx.h"
