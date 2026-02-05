#!/bin/sh
set -ex

ARCH=$1
SRC=test/e2e/testdata/main.c
mkdir -p binaries

. test/e2e/testdata/log-env.sh

build_c() { $1 $2 -o binaries/${ARCH}-$1-$3 $SRC; }
build_c_strip() { $1 $2 -o binaries/${ARCH}-$1-$3 $SRC && strip binaries/${ARCH}-$1-$3; }

build_c gcc "-Wl,-z,noexecstack" nx-explicit
build_c gcc "-Wl,-z,execstack" no-nx
build_c_strip gcc "-Wl,-z,noexecstack" nx-stripped
gcc -Wl,-z,noexecstack -static -o binaries/${ARCH}-gcc-nx-static $SRC
gcc -c -o binaries/${ARCH}-gcc-relocatable.o $SRC

build_c clang "-Wl,-z,noexecstack" nx-explicit
build_c clang "-Wl,-z,execstack" no-nx
build_c_strip clang "-Wl,-z,noexecstack" nx-stripped
clang -c -o binaries/${ARCH}-clang-relocatable.o $SRC

ls -la binaries/
