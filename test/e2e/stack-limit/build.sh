#!/bin/sh
set -ex

ARCH=$1
SRC=test/e2e/testdata/main.c
mkdir -p binaries

. test/e2e/testdata/log-env.sh

gcc -Wl,-z,stack-size=8388608 -o binaries/${ARCH}-gcc-stack-limit $SRC
gcc -o binaries/${ARCH}-gcc-no-stack-limit $SRC

clang -Wl,-z,stack-size=8388608 -o binaries/${ARCH}-clang-stack-limit $SRC
clang -o binaries/${ARCH}-clang-no-stack-limit $SRC

ls -la binaries/
