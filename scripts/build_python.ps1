# scripts/build_python.ps1
$outDir = "bindings/c/build"
$pkgDir = "python/gopptx"

if (!(Test-Path $outDir)) {
    New-Item -ItemType Directory -Force -Path $outDir
}

$libName = "gopptx.dll"
$libPath = Join-Path $outDir $libName
$releaseBuild = $env:GOPPTX_RELEASE_BUILD
$isReleaseBuild = $releaseBuild -eq "1" -or $releaseBuild -ieq "true"

Write-Host "Building Go engine for Python..."
if ($isReleaseBuild) {
    go build -trimpath -buildvcs=false -ldflags "-s -w" -o $libPath -buildmode=c-shared bindings/c/bridge.go
} else {
    go build -o $libPath -buildmode=c-shared bindings/c/bridge.go
}

if ($LASTEXITCODE -eq 0) {
    Write-Host "Build successful!"
    Write-Host "Copying $libName to Python package directory..."
    Copy-Item $libPath -Destination $pkgDir -Force
    Write-Host "Done! You can now install the package using 'pip install -e .'"
} else {
    Write-Error "Build failed!"
}
