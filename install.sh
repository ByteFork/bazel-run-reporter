#!/usr/bin/env sh
# bazel-run-reporter installer.
#
# Usage:
#   curl -fsSL https://install.bytefork.io/bazel-run-reporter | sh
#   curl -fsSL https://install.bytefork.io/bazel-run-reporter | sh -s -- --version v0.0.1 --bindir ~/.local/bin
#
# Alternate (direct from GitHub, no vanity redirect):
#   curl -fsSL https://raw.githubusercontent.com/ByteFork/bazel-run-reporter/main/install.sh | sh

set -eu

REPO="ByteFork/bazel-run-reporter"
BIN="bazel-run-reporter"
VERSION="latest"
BINDIR="/usr/local/bin"

while [ $# -gt 0 ]; do
  case "$1" in
    --version) VERSION="$2"; shift 2 ;;
    --bindir)  BINDIR="$2";  shift 2 ;;
    -h|--help)
      echo "Usage: install.sh [--version vX.Y.Z] [--bindir DIR]"
      echo "  --version  release tag (default: latest)"
      echo "  --bindir   install directory (default: /usr/local/bin)"
      exit 0
      ;;
    *) echo "unknown flag: $1" >&2; exit 1 ;;
  esac
done

# Resolve OS
OS=$(uname -s | tr '[:upper:]' '[:lower:]')
case "$OS" in
  darwin|linux) ;;
  *) echo "unsupported OS: $OS" >&2; exit 1 ;;
esac

# Resolve arch
ARCH=$(uname -m)
case "$ARCH" in
  x86_64|amd64) ARCH=amd64 ;;
  aarch64|arm64) ARCH=arm64 ;;
  *) echo "unsupported arch: $ARCH" >&2; exit 1 ;;
esac

# Resolve version
if [ "$VERSION" = "latest" ]; then
  VERSION=$(curl -fsSL "https://api.github.com/repos/${REPO}/releases/latest" \
    | sed -n 's/.*"tag_name": *"\([^"]*\)".*/\1/p' | head -1)
  [ -n "$VERSION" ] || { echo "failed to resolve latest version" >&2; exit 1; }
fi

# Strip leading v for asset filename; goreleaser omits it by default.
VERSION_NUM=${VERSION#v}

TMP=$(mktemp -d)
trap 'rm -rf "$TMP"' EXIT

ARCHIVE="${BIN}_${VERSION_NUM}_${OS}_${ARCH}.tar.gz"
URL="https://github.com/${REPO}/releases/download/${VERSION}/${ARCHIVE}"

echo "Downloading ${ARCHIVE}"
curl -fsSL "$URL" -o "$TMP/$ARCHIVE"

# Verify checksum against the checksums.txt from the release.
SUMS="${BIN}_${VERSION_NUM}_checksums.txt"
echo "Verifying SHA-256 against ${SUMS}"
curl -fsSL "https://github.com/${REPO}/releases/download/${VERSION}/${SUMS}" -o "$TMP/$SUMS"
(cd "$TMP" && grep " $ARCHIVE\$" "$SUMS" | shasum -a 256 -c -)

tar -xzf "$TMP/$ARCHIVE" -C "$TMP"

# Install
if [ ! -w "$BINDIR" ]; then
  echo "Installing to $BINDIR (requires sudo)"
  sudo install -m 0755 "$TMP/$BIN" "$BINDIR/$BIN"
else
  install -m 0755 "$TMP/$BIN" "$BINDIR/$BIN"
fi

echo "Installed at: $BINDIR/$BIN"
