#!/bin/sh
set -ex

SRC=test/e2e/testdata/main.c
mkdir -p binaries

. test/e2e/testdata/log-env.sh

ARCH=$(uname -m)
if [ "$ARCH" != "aarch64" ]; then
    echo "Error: ARM PAC is only supported on aarch64, detected $ARCH"
    exit 1
fi

build_c() { $1 $2 -o binaries/$1-$3 $SRC; }
build_c_strip() { $1 $2 -o binaries/$1-$3 $SRC && strip binaries/$1-$3; }

build_c gcc "-mbranch-protection=pac-ret" pac-enabled
build_c gcc "-mbranch-protection=none" pac-disabled
build_c_strip gcc "-mbranch-protection=pac-ret" pac-stripped

build_c clang "-mbranch-protection=pac-ret" pac-enabled
build_c clang "-mbranch-protection=none" pac-disabled
build_c_strip clang "-mbranch-protection=pac-ret" pac-stripped

ls -la binaries/
