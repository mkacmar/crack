#!/bin/sh
set -ex

ARCH=$1
SRC=test/e2e/testdata/main.c
mkdir -p binaries

. test/e2e/testdata/log-env.sh

gcc -fsanitize=address -o binaries/${ARCH}-gcc-asan $SRC
gcc -fsanitize=address -o binaries/${ARCH}-gcc-asan-stripped $SRC
strip binaries/${ARCH}-gcc-asan-stripped

clang -fsanitize=address -o binaries/${ARCH}-clang-asan $SRC
clang -fsanitize=address -o binaries/${ARCH}-clang-asan-stripped $SRC
strip binaries/${ARCH}-clang-asan-stripped

gcc -o binaries/${ARCH}-gcc-no-asan $SRC
clang -o binaries/${ARCH}-clang-no-asan $SRC

ls -la binaries/
