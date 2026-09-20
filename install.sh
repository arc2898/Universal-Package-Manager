#!/usr/bin/env bash

# UPM Universal Package Manager - Enhanced Installer
# Installs dependencies, builds, tests, and installs UPM binary

set -euo pipefail

ROOT_DIR="$(cd -- "$(dirname -- "${BASH_SOURCE[0]}")/.." && pwd)"
BINARY="${ROOT_DIR}/upm-bin"
INSTALL_DIR="/usr/local/bin"
LOG_DIR="/var/log"
LOG_FILE="${LOG_DIR}/upm.log"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

log_info() { echo -e "${BLUE}[INFO]${NC} $*"; }
log_success() { echo -e "${GREEN}[OK]${NC} $*"; }
log_warn() { echo -e "${YELLOW}[WARN]${NC} $*"; }
log_error() { echo -e "${RED}[ERROR]${NC} $*"; }

# Detect OS
detect_os() {
    if [[ "$OSTYPE" == "linux-gnu"* ]]; then
        if [[ -f /etc/os-release ]]; then
            . /etc/os-release
            OS="$ID"
            OS_LIKE="${ID_LIKE:-}"
        else
            OS="linux"
        fi
    elif [[ "$OSTYPE" == "darwin"* ]]; then
        OS="macos"
    elif [[ "$OSTYPE" == "msys" ]] || [[ "$OSTYPE" == "cygwin" ]] || [[ "$OSTYPE" == "win32" ]]; then
        OS="windows"
    else
        OS="unknown"
    fi
    echo "$OS"
}

OS=$(detect_os)
log_info "Detected OS: $OS"

# Check if running as root for system install
check_root() {
    if [[ $EUID -ne 0 ]] && [[ "$OS" != "windows" ]] && [[ "$USER_INSTALL" != "true" ]]; then
        log_warn "Not running as root. Will use sudo for system installation."
        SUDO="sudo"
    else
        SUDO=""
    fi
}

# Install Go if not present
install_go() {
    if command -v go &> /dev/null; then
        log_success "Go already installed: $(go version)"
        return 0
    fi

    log_info "Installing Go..."
    case "$OS" in
        ubuntu|debian|linuxmint|pop|elementary)
            $SUDO apt-get update && $SUDO apt-get install -y golang-go
            ;;
        fedora|rhel|centos|rocky|almalinux)
            $SUDO dnf install -y golang
            ;;
        arch|manjaro|endeavouros|garuda)
            $SUDO pacman -S --noconfirm go
            ;;
        opensuse*|suse|sles)
            $SUDO zypper install -y go
            ;;
        void)
            $SUDO xbps-install -S go
            ;;
        alpine)
            $SUDO apk add go
            ;;
        gentoo|calculate)
            $SUDO emerge --ask dev-lang/go
            ;;
        macos)
            if command -v brew &> /dev/null; then
                brew install go
            else
                log_error "Homebrew not found. Please install Go manually from https://golang.org/dl/"
                exit 1
            fi
            ;;
        windows)
            log_error "On Windows, please install Go manually from https://golang.org/dl/"
            log_error "Or use winget: winget install GoLang.Go"
            exit 1
            ;;
        *)
            log_error "Unsupported OS for automatic Go installation: $OS"
            log_error "Please install Go manually from https://golang.org/dl/"
            exit 1
            ;;
    esac
    log_success "Go installed successfully"
}

# Install build dependencies
install_build_deps() {
    log_info "Installing build dependencies..."
    case "$OS" in
        ubuntu|debian|linuxmint|pop|elementary)
            $SUDO apt-get update && $SUDO apt-get install -y build-essential git make
            ;;
        fedora|rhel|centos|rocky|almalinux)
            $SUDO dnf groupinstall -y "Development Tools" && $SUDO dnf install -y git make
            ;;
        arch|manjaro|endeavouros|garuda)
            $SUDO pacman -S --noconfirm base-devel git make
            ;;
        opensuse*|suse|sles)
            $SUDO zypper install -y -t pattern devel_basis && $SUDO zypper install -y git make
            ;;
        void)
            $SUDO xbps-install -S base-devel git make
            ;;
        alpine)
            $SUDO apk add build-base git make
            ;;
        gentoo|calculate)
            $SUDO emerge --ask sys-devel/gcc sys-devel/make dev-vcs/git
            ;;
        macos)
            if command -v brew &> /dev/null; then
                brew install make git
            else
                xcode-select --install 2>/dev/null || true
            fi
            ;;
        windows)
            log_info "Windows: build tools should be available via Go toolchain"
            ;;
        *)
            log_warn "Unknown OS, skipping build dependency installation"
            ;;
    esac
    log_success "Build dependencies installed"
}

# Build UPM
build_upm() {
    log_info "Building UPM..."
    cd "$ROOT_DIR"

    # Download dependencies
    log_info "Downloading Go modules..."
    go mod download

    # Run tests first
    log_info "Running tests..."
    if go test ./... -v; then
        log_success "All tests passed"
    else
        log_error "Tests failed! Build aborted."
        exit 1
    fi

    # Build binary
    log_info "Compiling binary..."
    if go build -ldflags="-s -w" -o "$BINARY" .; then
        log_success "Binary built at $BINARY"
    else
        log_error "Build failed!"
        exit 1
    fi

    # Verify binary works
    if "$BINARY" -v &> /dev/null; then
        log_success "Binary verification passed"
    else
        log_error "Binary verification failed!"
        exit 1
    fi
}

# Test UPM functionality
test_upm_functionality() {
    log_info "Testing UPM functionality..."

    # Test version flag
    log_info "Testing: upm -v"
    "$BINARY" -v

    # Test help
    log_info "Testing: upm -h"
    "$BINARY" -h 2>&1 | head -20

    # Test manager detection (non-destructive)
    log_info "Testing: upm refresh"
    "$BINARY" refresh 2>&1 | head -30

    log_success "Functionality tests passed"
}

# Install binary to system
install_binary() {
    log_info "Installing UPM to $INSTALL_DIR..."
    $SUDO install -m 0755 "$BINARY" "$INSTALL_DIR/upm"
    log_success "UPM installed to $INSTALL_DIR/upm"
}

# Setup logging
setup_logging() {
    log_info "Setting up logging..."
    $SUDO install -m 0640 /dev/null "$LOG_FILE" 2>/dev/null || {
        $SUDO touch "$LOG_FILE"
        $SUDO chmod 640 "$LOG_FILE"
    }
    $SUDO chown root:adm "$LOG_FILE" 2>/dev/null || true
    log_success "Log file created at $LOG_FILE"
}

# Verify installation
verify_installation() {
    log_info "Verifying installation..."
    # Add install dir to PATH for verification
    export PATH="$INSTALL_DIR:$PATH"
    if command -v upm &> /dev/null; then
        log_success "UPM is available in PATH"
        upm -v
        upm refresh 2>&1 | head -20
    else
        log_error "UPM not found in PATH after installation!"
        exit 1
    fi
}

# Cleanup build artifacts
cleanup() {
    log_info "Cleaning up build artifacts..."
    rm -f "$BINARY"
    log_success "Cleanup complete"
}

# Print usage
usage() {
    cat <<EOF
UPM Universal Package Manager - Installer

Usage: $0 [options]

Options:
    --no-test       Skip running tests before install
    --no-build-deps Skip build dependency installation
    --user          Install to ~/.local/bin instead of /usr/local/bin
    --help          Show this help

Examples:
    $0                    # Full install with tests
    $0 --user             # Install to user directory (no sudo needed)
    $0 --no-test          # Skip tests (faster)
EOF
}

# Parse arguments
SKIP_TESTS=false
SKIP_BUILD_DEPS=false
USER_INSTALL=false

for arg in "$@"; do
    case $arg in
        --no-test)
            SKIP_TESTS=true
            shift
            ;;
        --no-build-deps)
            SKIP_BUILD_DEPS=true
            shift
            ;;
        --user)
            USER_INSTALL=true
            INSTALL_DIR="$HOME/.local/bin"
            SUDO=""
            LOG_DIR="$HOME/.local/log"
            LOG_FILE="$LOG_DIR/upm.log"
            mkdir -p "$INSTALL_DIR" "$LOG_DIR"
            shift
            ;;
        --help|-h)
            usage
            exit 0
            ;;
        *)
            log_error "Unknown option: $arg"
            usage
            exit 1
            ;;
    esac
done

# Main installation flow
main() {
    log_info "=== UPM Universal Package Manager Installer ==="
    log_info "Root directory: $ROOT_DIR"
    log_info "Install directory: $INSTALL_DIR"
    log_info "Log file: $LOG_FILE"

    check_root
    install_go

    if [[ "$SKIP_BUILD_DEPS" != "true" ]]; then
        install_build_deps
    fi

    build_upm

    if [[ "$SKIP_TESTS" != "true" ]]; then
        test_upm_functionality
    fi

    install_binary
    setup_logging
    verify_installation
    cleanup

    log_success "=== Installation Complete ==="
    echo
    echo "Try running:"
    echo "  upm -h              # Show help"
    echo "  upm refresh         # Detect package managers"
    echo "  upm install git     # Install a package"
    echo "  upm -b install apt git,vim snap spotify  # Bulk install"
    echo "  upm -u              # Update all managers"
    echo "  upm -l              # View logs"
}

main "$@"