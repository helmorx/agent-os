$ErrorActionPreference = "Stop"

$Repo = $env:HELMOR_REPO
if ([string]::IsNullOrWhiteSpace($Repo)) { $Repo = "helmorx/agent-os" }

$Version = $env:HELMOR_VERSION
if ([string]::IsNullOrWhiteSpace($Version)) { $Version = "latest" }

$InstallDir = $env:HELMOR_INSTALL_DIR
if ([string]::IsNullOrWhiteSpace($InstallDir)) { $InstallDir = Join-Path $HOME ".helmor\bin" }

$Arch = if ([System.Runtime.InteropServices.RuntimeInformation]::OSArchitecture -eq "Arm64") { "arm64" } else { "amd64" }
$Asset = "helmor_windows_$Arch.zip"

if ($Version -eq "latest") {
  $Base = "https://github.com/$Repo/releases/latest/download"
} else {
  $Base = "https://github.com/$Repo/releases/download/$Version"
}

$Temp = Join-Path ([System.IO.Path]::GetTempPath()) ("helmor-" + [System.Guid]::NewGuid().ToString())
New-Item -ItemType Directory -Path $Temp | Out-Null

try {
  $Zip = Join-Path $Temp $Asset
  $Checksums = Join-Path $Temp "checksums.txt"
  Invoke-WebRequest -Uri "$Base/$Asset" -OutFile $Zip
  Invoke-WebRequest -Uri "$Base/checksums.txt" -OutFile $Checksums

  $Expected = (Select-String -Path $Checksums -Pattern $Asset).Line.Split(" ")[0]
  $Actual = (Get-FileHash -Algorithm SHA256 $Zip).Hash.ToLowerInvariant()
  if ($Expected.ToLowerInvariant() -ne $Actual) {
    throw "checksum mismatch for $Asset"
  }

  Expand-Archive -Path $Zip -DestinationPath $Temp -Force
  New-Item -ItemType Directory -Path $InstallDir -Force | Out-Null
  Copy-Item -Path (Join-Path $Temp "helmor.exe") -Destination (Join-Path $InstallDir "helmor.exe") -Force
  Write-Host "Installed helmor to $InstallDir\helmor.exe"
  Write-Host "Add $InstallDir to PATH if it is not already present."
} finally {
  Remove-Item -Recurse -Force $Temp
}

