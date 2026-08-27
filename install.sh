#!/bin/sh
set -eu

REPO="mkacmar/crack"
BINARY="crack"
SUMS="SHA256SUMS"

VERSION=""
INSTALL_DIR="$HOME/.local/bin"

usage() {
    echo "Usage: install.sh [--version <tag>] [--dir <path>]"
    echo ""
    echo "Download, verify and install the crack binary."
    echo ""
    echo "Options:"
    echo "  --version <tag>  Release to install (default: the latest release)"
    echo "  --dir <path>     Install directory (default: \$HOME/.local/bin)"
    echo ""
    echo "The download is checked against $SUMS. See the README for how to also verify the signature over it."
}

fail() {
    echo "Error: $*" >&2
    exit 1
}

usage_error() {
    [ $# -eq 0 ] || echo "Error: $*" >&2
    usage >&2
    exit 1
}

while [ $# -gt 0 ]; do
    case "$1" in
        --version)
            [ $# -ge 2 ] || usage_error "--version needs a value"
            VERSION="$2"
            shift 2
            ;;
        --dir)
            [ $# -ge 2 ] || usage_error "--dir needs a value"
            INSTALL_DIR="$2"
            shift 2
            ;;
        -h|--help)
            usage
            exit 0
            ;;
        *)
            usage_error "unknown option $1"
            ;;
    esac
done

case "$(uname -s)" in
    Linux)  OS="linux" ;;
    Darwin) OS="darwin" ;;
    *)      fail "this script supports Linux and macOS. Windows builds are on the releases page" ;;
esac

case "$(uname -m)" in
    x86_64|amd64)  ARCH="amd64" ;;
    aarch64|arm64) ARCH="arm64" ;;
    *)             fail "no prebuilt binary for $(uname -m), build from source instead" ;;
esac

ASSET="${BINARY}_${OS}_${ARCH}"

if [ -n "$VERSION" ]; then
    BASE_URL="https://github.com/$REPO/releases/download/$VERSION"
else
    BASE_URL="https://github.com/$REPO/releases/latest/download"
fi

if command -v curl >/dev/null 2>&1; then
    download() { curl -fsSL -o "$2" "$1"; }
elif command -v wget >/dev/null 2>&1; then
    download() { wget -q -O "$2" "$1"; }
else
    fail "neither curl nor wget is available"
fi

if command -v sha256sum >/dev/null 2>&1; then
    sha256_of() { sha256sum "$1" | cut -d ' ' -f 1; }
elif command -v shasum >/dev/null 2>&1; then
    sha256_of() { shasum -a 256 "$1" | cut -d ' ' -f 1; }
else
    fail "neither sha256sum nor shasum is available, cannot verify the download"
fi

WORK_DIR=$(mktemp -d)
STAGED="$INSTALL_DIR/.$BINARY.incomplete"
trap 'rm -rf "$WORK_DIR" "$STAGED"' EXIT INT TERM

echo "Downloading $ASSET from $BASE_URL"
download "$BASE_URL/$ASSET" "$WORK_DIR/$ASSET" || fail "could not download $BASE_URL/$ASSET"
download "$BASE_URL/$SUMS" "$WORK_DIR/$SUMS" || fail "could not download $BASE_URL/$SUMS"

EXPECTED=$(awk -v asset="$ASSET" '$2 == asset { print $1 }' "$WORK_DIR/$SUMS")
[ -n "$EXPECTED" ] || fail "$SUMS lists no entry for $ASSET"

ACTUAL=$(sha256_of "$WORK_DIR/$ASSET")
[ "$ACTUAL" = "$EXPECTED" ] || fail "checksum mismatch for $ASSET, refusing to install"
echo "Checksum verified"

mkdir -p "$INSTALL_DIR"
cp "$WORK_DIR/$ASSET" "$STAGED"
chmod 755 "$STAGED"
mv "$STAGED" "$INSTALL_DIR/$BINARY"

echo "Installed $("$INSTALL_DIR/$BINARY" --version | head -1) to $INSTALL_DIR/$BINARY"

case ":$PATH:" in
    *":$INSTALL_DIR:"*) ;;
    *) echo "Note: $INSTALL_DIR is not on your PATH" ;;
esac
