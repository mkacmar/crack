#!/bin/sh
set -ex

ARCH=$1
SRC=test/e2e/testdata/main.c
mkdir -p binaries

. test/e2e/testdata/log-env.sh

gcc -Wl,-z,separate-code -o binaries/${ARCH}-gcc-separate-code $SRC
gcc -Wl,-z,separate-code -o binaries/${ARCH}-gcc-separate-code-stripped $SRC
strip binaries/${ARCH}-gcc-separate-code-stripped
gcc -Wl,-z,separate-code -static -o binaries/${ARCH}-gcc-separate-code-static $SRC || echo "static linking not supported"
gcc -Wl,-z,separate-code -shared -fPIC -o binaries/${ARCH}-gcc-separate-code-shared $SRC

clang -Wl,-z,separate-code -o binaries/${ARCH}-clang-separate-code $SRC
clang -Wl,-z,separate-code -o binaries/${ARCH}-clang-separate-code-stripped $SRC
strip binaries/${ARCH}-clang-separate-code-stripped

ls -la binaries/

