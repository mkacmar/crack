#!/bin/sh
set -ex

ARCH=$1
SRC=test/e2e/testdata/main.c
mkdir -p binaries

. test/e2e/testdata/log-env.sh

gcc -Wl,-z,nodump -o binaries/${ARCH}-gcc-nodump $SRC
gcc -Wl,-z,nodump -o binaries/${ARCH}-gcc-nodump-stripped $SRC
strip binaries/${ARCH}-gcc-nodump-stripped

clang -Wl,-z,nodump -o binaries/${ARCH}-clang-nodump $SRC
clang -Wl,-z,nodump -o binaries/${ARCH}-clang-nodump-stripped $SRC
strip binaries/${ARCH}-clang-nodump-stripped

gcc -o binaries/${ARCH}-gcc-default $SRC
clang -o binaries/${ARCH}-clang-default $SRC

ls -la binaries/
