#!/bin/sh
set -ex

LIBC=$1
if [ -z "$LIBC" ]; then
    echo "Usage: $0 <glibc|musl>"
    exit 1
fi

SRC=test/e2e/testdata/main.c
mkdir -p binaries

. test/e2e/testdata/log-env.sh

ARCH=$(uname -m)
if [ "$ARCH" != "aarch64" ]; then
    echo "Error: ARM BTI is only supported on aarch64, detected $ARCH"
    exit 1
fi

if [ "$LIBC" = "musl" ]; then
    PREFIX="musl-"
else
    PREFIX=""
fi

build_c() { $1 $2 -o binaries/${PREFIX}$1-$3 $SRC; }
build_c_strip() { $1 $2 -o binaries/${PREFIX}$1-$3 $SRC && strip binaries/${PREFIX}$1-$3; }

build_c gcc "-mbranch-protection=bti -Wl,-z,force-bti" bti-enabled
build_c gcc "-mbranch-protection=none" bti-disabled
build_c_strip gcc "-mbranch-protection=bti -Wl,-z,force-bti" bti-stripped

build_c clang "-mbranch-protection=bti -Wl,-z,force-bti" bti-enabled
build_c clang "-mbranch-protection=none" bti-disabled
build_c_strip clang "-mbranch-protection=bti -Wl,-z,force-bti" bti-stripped

ls -la binaries/
