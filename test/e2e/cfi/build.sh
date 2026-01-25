#!/bin/sh
set -ex

ARCH=$1
SRC=test/e2e/testdata/cfi.c
mkdir -p binaries

. test/e2e/testdata/log-env.sh

# Cross-DSO CFI with LTO
clang -fsanitize=cfi -fsanitize-cfi-cross-dso -flto -fvisibility=hidden -fuse-ld=lld -o binaries/${ARCH}-clang-cfi $SRC
clang -fsanitize=cfi -fsanitize-cfi-cross-dso -flto -fvisibility=hidden -fuse-ld=lld -o binaries/${ARCH}-clang-cfi-stripped $SRC
strip binaries/${ARCH}-clang-cfi-stripped

# No CFI
clang -o binaries/${ARCH}-clang-no-cfi $SRC
clang -o binaries/${ARCH}-clang-no-cfi-stripped $SRC
strip binaries/${ARCH}-clang-no-cfi-stripped

gcc -o binaries/${ARCH}-gcc-no-cfi $SRC
gcc -o binaries/${ARCH}-gcc-no-cfi-stripped $SRC
strip binaries/${ARCH}-gcc-no-cfi-stripped

ls -la binaries/
