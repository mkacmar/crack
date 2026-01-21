#!/bin/sh
set -ex

SRC=test/e2e/testdata/main.c
mkdir -p binaries

. test/e2e/testdata/log-env.sh

ARCH=$(uname -m)
if [ "$ARCH" != "x86_64" ]; then
    echo "Error: Intel CET is only supported on x86_64, detected $ARCH"
    exit 1
fi

gcc -fcf-protection=full -o binaries/gcc-cet-full $SRC
gcc -fcf-protection=branch -o binaries/gcc-cet-branch $SRC
gcc -fcf-protection=none -o binaries/gcc-cet-none $SRC
gcc -fcf-protection=full -o binaries/gcc-cet-full-stripped $SRC
strip binaries/gcc-cet-full-stripped

clang -fcf-protection=full -o binaries/clang-cet-full $SRC
clang -fcf-protection=branch -o binaries/clang-cet-branch $SRC
clang -fcf-protection=none -o binaries/clang-cet-none $SRC

ls -la binaries/

