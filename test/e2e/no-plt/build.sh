#!/bin/sh
set -ex

ARCH=$1
SRC=test/e2e/testdata/plt.c
mkdir -p binaries

. test/e2e/testdata/log-env.sh

gcc -fno-plt -o binaries/${ARCH}-gcc-no-plt $SRC
gcc -fno-plt -o binaries/${ARCH}-gcc-no-plt-stripped $SRC
strip binaries/${ARCH}-gcc-no-plt-stripped

clang -fno-plt -o binaries/${ARCH}-clang-no-plt $SRC
clang -fno-plt -o binaries/${ARCH}-clang-no-plt-stripped $SRC
strip binaries/${ARCH}-clang-no-plt-stripped

gcc -o binaries/${ARCH}-gcc-default $SRC
clang -o binaries/${ARCH}-clang-default $SRC

ls -la binaries/
