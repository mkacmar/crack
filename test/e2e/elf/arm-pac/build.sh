#!/bin/sh
set -ex

C_SRC=test/e2e/elf/testdata/main.c
RUST_SRC=test/e2e/elf/testdata/main.rs
mkdir -p binaries

. test/e2e/elf/testdata/log-env.sh

ARCH=$(uname -m)
if [ "$ARCH" != "arm64" ]; then
    echo "Error: ARM PAC is only supported on arm64, detected $ARCH"
    exit 1
fi

build_c() { $1 $2 -o binaries/$1-$3 $C_SRC; }
build_c_strip() { $1 $2 -o binaries/$1-$3 $C_SRC && strip binaries/$1-$3; }

build_c gcc "-mbranch-protection=pac-ret" pac-enabled
build_c gcc "-mbranch-protection=none" pac-disabled
build_c_strip gcc "-mbranch-protection=pac-ret" pac-stripped

build_c clang "-mbranch-protection=pac-ret" pac-enabled
build_c clang "-mbranch-protection=none" pac-disabled
build_c_strip clang "-mbranch-protection=pac-ret" pac-stripped

rustc -o binaries/rustc-no-pac $RUST_SRC

ls -la binaries/
