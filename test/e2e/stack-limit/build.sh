#!/bin/sh
set -ex

ARCH=$1
SRC=test/e2e/testdata/main.c
mkdir -p binaries

. test/e2e/testdata/log-env.sh

build_c() { $1 $2 -o binaries/${ARCH}-$1-$3 $SRC; }

build_c gcc "-Wl,-z,stack-size=8388608" stack-limit
build_c gcc "" no-stack-limit
gcc -c -o binaries/${ARCH}-gcc-relocatable.o $SRC

build_c clang "-Wl,-z,stack-size=8388608" stack-limit
build_c clang "" no-stack-limit
clang -c -o binaries/${ARCH}-clang-relocatable.o $SRC

ls -la binaries/
