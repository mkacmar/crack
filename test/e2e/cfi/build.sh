#!/bin/sh
set -ex

ARCH=$1
SRC=test/e2e/testdata/cfi.c
mkdir -p binaries

. test/e2e/testdata/log-env.sh

clang -fsanitize=cfi -flto -fvisibility=hidden -o binaries/${ARCH}-clang-cfi $SRC
clang -fsanitize=cfi -flto -fvisibility=hidden -o binaries/${ARCH}-clang-cfi-stripped $SRC
strip binaries/${ARCH}-clang-cfi-stripped

clang -o binaries/${ARCH}-clang-no-cfi $SRC
gcc -o binaries/${ARCH}-gcc-no-cfi $SRC

ls -la binaries/
