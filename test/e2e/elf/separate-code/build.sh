#!/bin/sh
set -ex

ARCH=$1
C_SRC=test/e2e/elf/testdata/main.c
RUST_SRC=test/e2e/elf/testdata/main.rs
mkdir -p binaries

. test/e2e/elf/testdata/log-env.sh

build_c() { $1 $2 -o binaries/${ARCH}-$1-$3 $C_SRC; }
build_c_strip() { $1 $2 -o binaries/${ARCH}-$1-$3 $C_SRC && strip binaries/${ARCH}-$1-$3; }

build_c gcc "-Wl,-z,separate-code" separate-code
build_c_strip gcc "-Wl,-z,separate-code" separate-code-stripped
gcc -Wl,-z,separate-code -static -o binaries/${ARCH}-gcc-separate-code-static $C_SRC
build_c gcc "-Wl,-z,separate-code -shared -fPIC" separate-code-shared

build_c clang "-Wl,-z,separate-code" separate-code
build_c_strip clang "-Wl,-z,separate-code" separate-code-stripped
clang -c -o binaries/${ARCH}-clang-relocatable.o $C_SRC

gcc -c -o binaries/${ARCH}-gcc-relocatable.o $C_SRC

rustc -C link-arg=-z -C link-arg=separate-code -o binaries/${ARCH}-rustc-separate-code $RUST_SRC

ls -la binaries/
