#!/bin/sh
set -ex

ARCH=$1
SRC=test/e2e/testdata/main.c
mkdir -p binaries

. test/e2e/testdata/log-env.sh

gcc -Wl,-z,noexecstack -o binaries/${ARCH}-gcc-nx-explicit $SRC
gcc -Wl,-z,execstack -o binaries/${ARCH}-gcc-no-nx $SRC
gcc -Wl,-z,noexecstack -o binaries/${ARCH}-gcc-nx-stripped $SRC
strip binaries/${ARCH}-gcc-nx-stripped
gcc -Wl,-z,noexecstack -static -o binaries/${ARCH}-gcc-nx-static $SRC || echo "static linking not supported"
gcc -c -o binaries/${ARCH}-gcc-relocatable.o $SRC

clang -Wl,-z,noexecstack -o binaries/${ARCH}-clang-nx-explicit $SRC
clang -Wl,-z,execstack -o binaries/${ARCH}-clang-no-nx $SRC
clang -Wl,-z,noexecstack -o binaries/${ARCH}-clang-nx-stripped $SRC
strip binaries/${ARCH}-clang-nx-stripped
clang -c -o binaries/${ARCH}-clang-relocatable.o $SRC

ls -la binaries/
