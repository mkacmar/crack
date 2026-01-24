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

gcc -mbranch-protection=pac-ret -o binaries/gcc-pac-enabled $SRC
gcc -mbranch-protection=none -o binaries/gcc-pac-disabled $SRC
gcc -mbranch-protection=pac-ret -o binaries/gcc-pac-stripped $SRC
strip binaries/gcc-pac-stripped

clang -mbranch-protection=pac-ret -o binaries/clang-pac-enabled $SRC
clang -mbranch-protection=none -o binaries/clang-pac-disabled $SRC
clang -mbranch-protection=pac-ret -o binaries/clang-pac-stripped $SRC
strip binaries/clang-pac-stripped

ls -la binaries/
