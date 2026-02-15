#!/bin/sh
set -ex

ARCH=$1
C_SRC=test/e2e/elf/testdata/main.c
RUST_SRC=test/e2e/elf/testdata/main.rs
mkdir -p binaries

. test/e2e/elf/testdata/log-env.sh

build_c() { $1 $2 -o binaries/${ARCH}-$1-$3 $C_SRC; }
build_c_strip() { $1 $2 -o binaries/${ARCH}-$1-$3 $C_SRC && strip binaries/${ARCH}-$1-$3; }

build_c gcc "-Wl,-z,nodump" nodump
build_c_strip gcc "-Wl,-z,nodump" nodump-stripped
gcc -c -o binaries/${ARCH}-gcc-relocatable.o $C_SRC

build_c clang "-Wl,-z,nodump" nodump
build_c_strip clang "-Wl,-z,nodump" nodump-stripped
clang -c -o binaries/${ARCH}-clang-relocatable.o $C_SRC

build_c gcc "" default
build_c clang "" default


rustc -o binaries/${ARCH}-rustc-default $RUST_SRC

ls -la binaries/
