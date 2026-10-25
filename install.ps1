# UPM Universal Package Manager - Windows PowerShell Installer

$ErrorActionPreference = "Stop"

# Colors for output
function Write-Log-Info($msg) { Write-Host "[INFO] $msg" -ForegroundColor Cyan }
function Write-Log-Success($msg) { Write-Host "[OK] $msg" -ForegroundColor Green }
function Write-Log-Warn($msg) { Write-Host "[WARN] $msg" -ForegroundColor Yellow }
function Write-Log-Error($msg) { Write-Host "[ERROR] $msg" -ForegroundColor Red }

Write-Log-Info "=== UPM Universal Package Manager Installer (Windows) ==="

# 1. Detect Go Installation
Write-Log-Info "Checking for Go installation..."
try {
    $goVer = go version
    Write-Log-Success "Go already installed: $goVer"
} catch {
    Write-Log-Error "Go is not installed."
    Write-Log-Warn "Please install Go using winget: winget install GoLang.Go"
    Write-Log-Warn "Or download it from: https://golang.org/dl/"
    exit 1
}

# 2. Build UPM
Write-Log-Info "Building UPM..."
try {
    Write-Log-Info "Downloading Go modules..."
    go mod download

    Write-Log-Info "Compiling binary..."
    # -s -w strips symbol table and debug info to reduce binary size
    go build -ldflags="-s -w" -o upm.exe .
    Write-Log-Success "Binary built as upm.exe"
} catch {
    Write-Log-Error "Build failed: $($_.Exception.Message)"
    exit 1
}

# 3. Run Tests
Write-Log-Info "Running tests..."
try {
    go test ./... -v
    Write-Log-Success "All tests passed"
} catch {
    Write-Log-Error "Tests failed! Installation aborted to ensure stability."
    exit 1
}

# 4. Install Binary
$InstallDir = "$HOME\upm"
if (-not (Test-Path $InstallDir)) {
    New-Item -ItemType Directory -Path $InstallDir | Out-Null
}

Write-Log-Info "Installing UPM to $InstallDir..."
Move-Item -Path "upm.exe" -Destination "$InstallDir\upm.exe" -Force

# 5. Update PATH Environment Variable
Write-Log-Info "Updating User PATH..."
$CurrentPath = [Environment]::GetEnvironmentVariable("Path", "User")
if ($CurrentPath -notlike "*$InstallDir*") {
    $NewPath = "$CurrentPath;$InstallDir"
    [Environment]::SetEnvironmentVariable("Path", $NewPath, "User")
    Write-Log-Success "Added $InstallDir to User PATH."
    Write-Log-Warn "NOTE: You may need to restart your terminal for the PATH change to take effect."
} else {
    Write-Log-Info "UPM directory already in PATH."
}

# 6. Verification
Write-Log-Info "Verifying installation..."
try {
    # Use absolute path for immediate verification before shell restart
    & "$InstallDir\upm.exe" -v
    Write-Log-Success "UPM verification passed!"
} catch {
    Write-Log-Error "Verification failed. Please check if the binary is compatible with your architecture."
    exit 1
}

Write-Log-Success "=== Installation Complete ==="
Write-Host "`nTry running (after restarting terminal):"
Write-Host "  upm -h              # Show help"
Write-Host "  upm refresh         # Detect package managers"
Write-Host "  upm install git     # Install a package"
