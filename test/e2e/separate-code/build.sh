#!/bin/sh
set -ex

ARCH=$1
SRC=test/e2e/testdata/main.c
mkdir -p binaries

. test/e2e/testdata/log-env.sh

build() { $1 $2 -o binaries/${ARCH}-$1-$3 $SRC; }
build_strip() { $1 $2 -o binaries/${ARCH}-$1-$3 $SRC && strip binaries/${ARCH}-$1-$3; }

build gcc "-Wl,-z,separate-code" separate-code
build_strip gcc "-Wl,-z,separate-code" separate-code-stripped
gcc -Wl,-z,separate-code -static -o binaries/${ARCH}-gcc-separate-code-static $SRC
build gcc "-Wl,-z,separate-code -shared -fPIC" separate-code-shared

build clang "-Wl,-z,separate-code" separate-code
build_strip clang "-Wl,-z,separate-code" separate-code-stripped
clang -c -o binaries/${ARCH}-clang-relocatable.o $SRC

gcc -c -o binaries/${ARCH}-gcc-relocatable.o $SRC

ls -la binaries/
