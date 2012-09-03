$scriptpath = $MyInvocation.MyCommand.Path
$dir = Split-Path $scriptpath | Split-Path

$goPath = $env:GOPATH

if (-not $goPath.contains($dir)) {
    $env:GOPATH = "$goPath;$dir"
}