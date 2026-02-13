# scripts/build_python.ps1
$outDir = "bindings/c/build"
$pkgDir = "python/gopptx"

if (!(Test-Path $outDir)) {
    New-Item -ItemType Directory -Force -Path $outDir
}

$libName = "gopptx.dll"
$libPath = Join-Path $outDir $libName

Write-Host "Building Go engine for Python..."
go build -o $libPath -buildmode=c-shared bindings/c/bridge.go

if ($LASTEXITCODE -eq 0) {
    Write-Host "Build successful!"
    Write-Host "Copying $libName to Python package directory..."
    Copy-Item $libPath -Destination $pkgDir -Force
    Write-Host "Done! You can now install the package using 'pip install -e .'"
} else {
    Write-Error "Build failed!"
}
