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

gcc -mbranch-protection=bti -Wl,-z,force-bti -o binaries/${PREFIX}gcc-bti-enabled $SRC
gcc -mbranch-protection=none -o binaries/${PREFIX}gcc-bti-disabled $SRC
gcc -mbranch-protection=bti -Wl,-z,force-bti -o binaries/${PREFIX}gcc-bti-stripped $SRC
strip binaries/${PREFIX}gcc-bti-stripped

clang -mbranch-protection=bti -Wl,-z,force-bti -o binaries/${PREFIX}clang-bti-enabled $SRC
clang -mbranch-protection=none -o binaries/${PREFIX}clang-bti-disabled $SRC
clang -mbranch-protection=bti -Wl,-z,force-bti -o binaries/${PREFIX}clang-bti-stripped $SRC
strip binaries/${PREFIX}clang-bti-stripped

ls -la binaries/
