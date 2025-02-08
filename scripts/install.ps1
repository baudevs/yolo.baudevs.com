param(
    [string]$Version = "0.1.0"
)

# Detect architecture
$arch = if ([Environment]::Is64BitOperatingSystem) {
    if ([System.Runtime.InteropServices.RuntimeInformation]::ProcessArchitecture -eq [System.Runtime.InteropServices.Architecture]::Arm64) {
        "arm64"
    } else {
        "amd64"
    }
} else {
    Write-Error "Unsupported architecture. Only 64-bit systems are supported."
    exit 1
}

$binaryName = "yolo"
$installDir = "$env:LOCALAPPDATA\Programs\YOLO"
$binPath = "$installDir\$binaryName.exe"

# Create temporary directory
$tmpDir = New-TemporaryFile | ForEach-Object { Remove-Item $_; New-Item -ItemType Directory $_ }
Push-Location $tmpDir

# Download the appropriate release
$releaseUrl = "https://github.com/baudevs/yolo-cli/releases/download/v${Version}/${binaryName}-${Version}-windows-${arch}.zip"
Write-Host "Downloading YOLO CLI v${Version} for windows/${arch}..."
Invoke-WebRequest -Uri $releaseUrl -OutFile release.zip

# Extract the archive
Expand-Archive -Path release.zip -DestinationPath .

# Create installation directory if it doesn't exist
if (-not (Test-Path $installDir)) {
    New-Item -ItemType Directory -Path $installDir | Out-Null
}

# Install the binary
Write-Host "Installing YOLO CLI to $installDir..."
Copy-Item "${binaryName}-windows-${arch}\${binaryName}.exe" -Destination $binPath

# Add to PATH if not already present
$userPath = [Environment]::GetEnvironmentVariable("Path", "User")
if ($userPath -notlike "*$installDir*") {
    [Environment]::SetEnvironmentVariable(
        "Path",
        "$userPath;$installDir",
        "User"
    )
    Write-Host "Added YOLO CLI to your PATH"
}

# Clean up
Pop-Location
Remove-Item -Recurse -Force $tmpDir

Write-Host "✓ YOLO CLI has been installed successfully!"
Write-Host "Run 'yolo --help' to get started." 