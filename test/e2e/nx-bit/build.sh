#!/bin/sh
set -ex

ARCH=$1
SRC=test/e2e/testdata/main.c
mkdir -p binaries

. test/e2e/testdata/log-env.sh

build() { $1 $2 -o binaries/${ARCH}-$1-$3 $SRC; }
build_strip() { $1 $2 -o binaries/${ARCH}-$1-$3 $SRC && strip binaries/${ARCH}-$1-$3; }

build gcc "-Wl,-z,noexecstack" nx-explicit
build gcc "-Wl,-z,execstack" no-nx
build_strip gcc "-Wl,-z,noexecstack" nx-stripped
gcc -Wl,-z,noexecstack -static -o binaries/${ARCH}-gcc-nx-static $SRC
gcc -c -o binaries/${ARCH}-gcc-relocatable.o $SRC

build clang "-Wl,-z,noexecstack" nx-explicit
build clang "-Wl,-z,execstack" no-nx
build_strip clang "-Wl,-z,noexecstack" nx-stripped
clang -c -o binaries/${ARCH}-clang-relocatable.o $SRC

ls -la binaries/
