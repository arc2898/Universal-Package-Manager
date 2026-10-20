# Universal Package Manager (UPM)

[![Go](https://img.shields.io/badge/Go-1.22.2-00ADD8?logo=go&logoColor=white)](https://go.dev/)
[![License](https://img.shields.io/badge/License-MIT-green)](LICENSE)

UPM is a Go command-line tool that provides a unified interface for common package-management operations through system and language-specific package-manager adapters. Linux is the primary target; macOS and Windows adapters are best-effort and should be validated on their native operating systems before production use.

[![Version](https://img.shields.io/badge/version-0.2.0-00ADD8)](https://github.com/arc2898/Universal-Package-Manager/releases)

## Features

UPM supports installation, removal, real package searches, bulk installation and removal, full updates of detected managers, re-detection of known managers, operation logs, and self-removal of the installed binary. Package-manager operations are executed with argument-based process APIs rather than shell interpolation.

### Supported Package Managers

**System Managers:**
- **APT** - Debian/Ubuntu packages
- **DNF** - Fedora/RHEL/CentOS/RPM-based systems
- **Pacman** - Arch Linux/Manjaro packages
- **RPM** - Local RPM-file install/removal and searches among installed packages; repository-wide updates are delegated to the distribution manager
- **XBPS** - Void Linux packages
- **Zypper** - openSUSE/SUSE packages
- **Snap** - Snap packages
- **Flatpak** - Flatpak applications and runtimes
- **Homebrew** - macOS packages
- **Winget** - Windows packages (winget)
- **Chocolatey** - Windows packages (choco)
- **Scoop** - Windows packages (scoop)

**Language-Specific Managers:**
- **Cargo** - Rust packages
- **NPM** - Node.js packages
- **Pip** - Python packages (search uses `pip index versions` for an exact package lookup)
- **Gem** - Ruby packages

### New Features

- **Fallback search**: If a package is not found in one manager, UPM tries other available managers automatically
- **Multiple platform adapters**: Includes adapters for Linux distributions, macOS tools, and Windows package managers; native-platform verification is still required
- **Operation logs**: Root-run commands log to `/var/log/upm.log`; non-root commands use `~/.local/log/upm.log` by default. Set `UPM_LOG_FILE` to override the path
- **Self-removal**: `upm -r` removes the installed UPM binary on Unix; on Windows, close UPM and remove the executable manually

## Installation

From the repository root, run:

```bash
./install.sh
```

Alternatively:

```bash
make install
```

The installation builds the root Go package, installs `/usr/local/bin/upm` with mode `0755`, creates `/var/log/upm.log` with mode `0640`, and removes the temporary build artifact. The install script requires Go and system-install privileges (normally `sudo`); `--user` installs into `~/.local/bin` and uses a per-user log path.

## Usage

```text
upm install <pkg>       Install using the preferred native manager
upm remove <pkg>        Remove using the preferred native manager
upm search <pkg>        Search all detected managers
upm -b install ...      Bulk install using explicit manager/package pairs
upm -b remove ...       Bulk remove using explicit manager/package pairs
upm -u                  Update all detected managers
upm refresh             Re-detect known managers
upm -l                  View operation logs
upm -r                  Remove the installed UPM binary
upm -v                  Show the version
upm -h                  Show help
```

### Examples

```bash
upm install neovim
upm search firefox
upm -u
upm -b install apt htop,curl snap code
upm -l
```

Bulk arguments must contain complete manager/package pairs. Package names are comma-separated and empty names are rejected:

```bash
upm -b install apt git,vim snap code
```

UPM may request administrator privileges for system-wide APT, DNF, Pacman, RPM, XBPS, Zypper, Snap, and Portage operations. Flatpak operations use the normal Flatpak command and therefore follow Flatpak's system/user behavior. `upm -u` skips RPM, Cargo, and Pip because they lack a safe, repository-wide update operation through the current adapters.

## Development

Build, test, and statically analyze the project with:

```bash
make build
make test
go vet ./...
```

The GitHub Actions workflow runs tests, `go vet`, and a build on Linux, macOS, and Windows. This validates compilation and tests on those runners, but it does not replace end-to-end testing of real package-manager operations on each distribution.

To add a manager, create an adapter in `internal/adapters/`, implement `manager.Manager`, register it in `internal/core/core.go`, and add tests for command construction and output parsing. Tests must not invoke real package managers or mutate the host system.

### Project Structure

```text
.
├── main.go                    # CLI entry point and argument validation
├── internal/
│   ├── adapters/              # APT, DNF, Pacman, Snap, Flatpak, RPM, Homebrew, Winget, Chocolatey, Scoop, Cargo, NPM, Pip, Gem, and additional managers
│   ├── core/                  # Manager selection and operation orchestration
│   ├── detectors/             # Availability detection
│   └── logger/                # Operation log handling
├── pkg/manager/               # Manager interface and search result types
├── install.sh                 # Build and installation script
├── .github/workflows/ci.yml   # Cross-platform test/build workflow
├── go.mod
├── Makefile
├── LICENSE
└── README.md
```

## Logging and safety

Operational failures are returned to the CLI, printed to stderr, logged when possible, and reflected in a non-zero exit status. The system log is `/var/log/upm.log`; it is created with restrictive permissions rather than being world-writable. UPM does not run package installation, removal, or update commands during tests.

## License

This project is licensed under the MIT License. See [LICENSE](LICENSE).

## Verification

Before submitting changes, run `make build`, `make test`, and `go vet ./...`. Tests should use mocked manager commands and must not modify the host package database. Manual testing of actual installs, removals, or system-wide updates must be done in disposable VMs/containers for each target operating system; those operations are intentionally not run by CI.

## Roadmap

- Improve package-search accuracy and consistent output across adapters.
- Expand mocked tests for command construction, parsing, and failure cases.
- Improve installer reliability and make platform support clearer.

The current implementation is written in Go; a possible future rewrite is not part of the current release.