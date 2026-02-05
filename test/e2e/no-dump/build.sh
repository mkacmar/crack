#!/bin/sh
set -ex

ARCH=$1
SRC=test/e2e/testdata/main.c
mkdir -p binaries

. test/e2e/testdata/log-env.sh

build_c() { $1 $2 -o binaries/${ARCH}-$1-$3 $SRC; }
build_c_strip() { $1 $2 -o binaries/${ARCH}-$1-$3 $SRC && strip binaries/${ARCH}-$1-$3; }

build_c gcc "-Wl,-z,nodump" nodump
build_c_strip gcc "-Wl,-z,nodump" nodump-stripped
gcc -c -o binaries/${ARCH}-gcc-relocatable.o $SRC

build_c clang "-Wl,-z,nodump" nodump
build_c_strip clang "-Wl,-z,nodump" nodump-stripped
clang -c -o binaries/${ARCH}-clang-relocatable.o $SRC

build_c gcc "" default
build_c clang "" default

ls -la binaries/
