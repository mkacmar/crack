#!/bin/sh
set -ex

ARCH=$1
SRC=test/e2e/testdata/main.c
mkdir -p binaries

. test/e2e/testdata/log-env.sh

build() { $1 $2 -o binaries/${ARCH}-$1-$3 $SRC; }

build gcc "-Wl,-z,stack-size=8388608" stack-limit
build gcc "" no-stack-limit
gcc -c -o binaries/${ARCH}-gcc-relocatable.o $SRC

build clang "-Wl,-z,stack-size=8388608" stack-limit
build clang "" no-stack-limit
clang -c -o binaries/${ARCH}-clang-relocatable.o $SRC

ls -la binaries/
