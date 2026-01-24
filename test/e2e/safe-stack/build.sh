#!/bin/sh
set -ex

ARCH=$1
SRC=test/e2e/testdata/main.c
mkdir -p binaries

. test/e2e/testdata/log-env.sh

clang -fsanitize=safe-stack -o binaries/${ARCH}-clang-safestack $SRC
clang -fsanitize=safe-stack -o binaries/${ARCH}-clang-safestack-stripped $SRC
strip binaries/${ARCH}-clang-safestack-stripped

clang -o binaries/${ARCH}-clang-no-safestack $SRC
gcc -o binaries/${ARCH}-gcc-no-safestack $SRC

ls -la binaries/
