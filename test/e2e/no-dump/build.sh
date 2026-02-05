#!/bin/sh
set -ex

ARCH=$1
SRC=test/e2e/testdata/main.c
mkdir -p binaries

. test/e2e/testdata/log-env.sh

build() { $1 $2 -o binaries/${ARCH}-$1-$3 $SRC; }
build_strip() { $1 $2 -o binaries/${ARCH}-$1-$3 $SRC && strip binaries/${ARCH}-$1-$3; }

build gcc "-Wl,-z,nodump" nodump
build_strip gcc "-Wl,-z,nodump" nodump-stripped
gcc -c -o binaries/${ARCH}-gcc-relocatable.o $SRC

build clang "-Wl,-z,nodump" nodump
build_strip clang "-Wl,-z,nodump" nodump-stripped
clang -c -o binaries/${ARCH}-clang-relocatable.o $SRC

build gcc "" default
build clang "" default

ls -la binaries/
