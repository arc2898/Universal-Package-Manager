#!/usr/bin/env bash

set -euo pipefail

ROOT_DIR="$(cd -- "$(dirname -- "${BASH_SOURCE[0]}")/.." && pwd)"
BINARY="${ROOT_DIR}/upm-bin"

cd "$ROOT_DIR"
echo "Building UPM..."
go build -o "$BINARY" .

echo "Installing UPM to /usr/local/bin..."
sudo install -m 0755 "$BINARY" /usr/local/bin/upm

 echo "Setting up logs..."
sudo install -m 0640 /dev/null /var/log/upm.log
rm -f "$BINARY"

echo "UPM installed successfully!"
echo "Try running 'upm -h' to get started."
