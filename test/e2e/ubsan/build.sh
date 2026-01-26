#!/bin/sh
set -ex

ARCH=$1
SRC=test/e2e/testdata/ubsan.c
mkdir -p binaries

. test/e2e/testdata/log-env.sh

gcc -fsanitize=undefined -o binaries/${ARCH}-gcc-ubsan $SRC
gcc -fsanitize=undefined -o binaries/${ARCH}-gcc-ubsan-stripped $SRC
strip binaries/${ARCH}-gcc-ubsan-stripped

clang -fsanitize=undefined -o binaries/${ARCH}-clang-ubsan $SRC
clang -fsanitize=undefined -o binaries/${ARCH}-clang-ubsan-stripped $SRC
strip binaries/${ARCH}-clang-ubsan-stripped

gcc -o binaries/${ARCH}-gcc-no-ubsan $SRC
clang -o binaries/${ARCH}-clang-no-ubsan $SRC

ls -la binaries/
