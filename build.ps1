$ErrorActionPreference = "Stop"
Set-Location $PSScriptRoot

if (-not (Select-String -Path "version.go" -Pattern 'const Version = "([^"]+)"' -Quiet)) {
    throw "Cannot read version from version.go"
}
$version = (Select-String -Path "version.go" -Pattern 'const Version = "([^"]+)"').Matches[0].Groups[1].Value
$outputName = "env-editor-$version"

Write-Host "Building $outputName ..."

Write-Host "-> elevhelper.exe"
& go build -ldflags "-H windowsgui -s -w" -o elevhelper.exe ./cmd/elevhelper

Write-Host "-> frontend"
Push-Location frontend
& npm run build
Pop-Location

Write-Host "-> wails ($outputName.exe)"
& wails build -o $outputName

$binPath = Join-Path "build" "bin" "$outputName.exe"
if (-not (Test-Path $binPath)) {
    throw "Build failed: $binPath not found"
}

Write-Host "Done: $binPath"
