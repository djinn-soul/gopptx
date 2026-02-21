# scripts/build_cbridge.ps1
$outDir = "bindings/c/build"
if (!(Test-Path $outDir)) {
    New-Item -ItemType Directory -Force -Path $outDir
}

$dllPath = Join-Path $outDir "gopptx.dll"
$headerPath = Join-Path $outDir "gopptx.h"
$pythonDllPath = "python/gopptx/gopptx.dll"

Write-Host "Building gopptx shared library..."
go build -o $dllPath -buildmode=c-shared bindings/c/bridge.go

if ($LASTEXITCODE -eq 0) {
    Write-Host "Build successful!"
    Write-Host "DLL: $dllPath"
    Write-Host "Header: $headerPath"
    Copy-Item $dllPath -Destination $pythonDllPath -Force
    Write-Host "Python package DLL synced: $pythonDllPath"
} else {
    Write-Error "Build failed!"
}
