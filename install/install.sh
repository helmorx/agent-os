#!/usr/bin/env sh
set -eu

REPO="${HELMOR_REPO:-helmorx/helmoragent}"
VERSION="${HELMOR_VERSION:-latest}"
INSTALL_DIR="${HELMOR_INSTALL_DIR:-$HOME/.local/bin}"
OS="$(uname -s | tr '[:upper:]' '[:lower:]')"
ARCH="$(uname -m)"

case "$OS" in
  darwin|linux) ;;
  *) echo "unsupported operating system: $OS" >&2; exit 1 ;;
esac

case "$ARCH" in
  arm64|aarch64) ARCH="arm64" ;;
  x86_64|amd64) ARCH="amd64" ;;
  *) echo "unsupported architecture: $ARCH" >&2; exit 1 ;;
esac

ASSET="helmor_${OS}_${ARCH}.tar.gz"

if [ "$VERSION" = "latest" ]; then
  URL="https://github.com/$REPO/releases/latest/download/$ASSET"
  SUMS="https://github.com/$REPO/releases/latest/download/checksums.txt"
else
  URL="https://github.com/$REPO/releases/download/$VERSION/$ASSET"
  SUMS="https://github.com/$REPO/releases/download/$VERSION/checksums.txt"
fi

TMP_DIR="$(mktemp -d)"
trap 'rm -rf "$TMP_DIR"' EXIT

curl -fsSL "$URL" -o "$TMP_DIR/$ASSET"
curl -fsSL "$SUMS" -o "$TMP_DIR/checksums.txt"

if command -v sha256sum >/dev/null 2>&1; then
  (cd "$TMP_DIR" && grep "$ASSET" checksums.txt | sha256sum -c -)
elif command -v shasum >/dev/null 2>&1; then
  (cd "$TMP_DIR" && grep "$ASSET" checksums.txt | shasum -a 256 -c -)
else
  echo "warning: no sha256 verifier found; skipping checksum verification" >&2
fi

mkdir -p "$INSTALL_DIR"
tar -xzf "$TMP_DIR/$ASSET" -C "$TMP_DIR"
install "$TMP_DIR/helmor" "$INSTALL_DIR/helmor"
echo "Installed helmor to $INSTALL_DIR/helmor"
