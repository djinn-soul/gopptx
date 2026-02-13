#!/bin/bash
# scripts/build_python.sh
set -e

outDir="bindings/c/build"
pkgDir="python/gopptx"
mkdir -p "$outDir"

libName="libgopptx.so"
if [[ "$OSTYPE" == "darwin"* ]]; then
    libName="libgopptx.dylib"
fi

echo "Building Go engine for Python..."
go build -o "$outDir/$libName" -buildmode=c-shared bindings/c/bridge.go

echo "Build successful!"
echo "Copying $libName to Python package directory..."
cp "$outDir/$libName" "$pkgDir/"

echo "Done! You can now install the package using 'pip install -e .'"
