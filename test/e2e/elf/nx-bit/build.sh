#!/bin/sh
set -ex

ARCH=$1
C_SRC=test/e2e/elf/testdata/main.c
RUST_SRC=test/e2e/elf/testdata/main.rs
mkdir -p binaries

. test/e2e/elf/testdata/log-env.sh

build_c() { $1 $2 -o binaries/${ARCH}-$1-$3 $C_SRC; }
build_c_strip() { $1 $2 -o binaries/${ARCH}-$1-$3 $C_SRC && strip binaries/${ARCH}-$1-$3; }

build_c gcc "-Wl,-z,noexecstack" nx-explicit
build_c gcc "-Wl,-z,execstack" no-nx
build_c_strip gcc "-Wl,-z,noexecstack" nx-stripped
gcc -Wl,-z,noexecstack -static -o binaries/${ARCH}-gcc-nx-static $C_SRC
gcc -c -o binaries/${ARCH}-gcc-relocatable.o $C_SRC

build_c clang "-Wl,-z,noexecstack" nx-explicit
build_c clang "-Wl,-z,execstack" no-nx
build_c_strip clang "-Wl,-z,noexecstack" nx-stripped
clang -c -o binaries/${ARCH}-clang-relocatable.o $C_SRC

rustc -o binaries/${ARCH}-rustc-nx $RUST_SRC
rustc -C strip=symbols -o binaries/${ARCH}-rustc-nx-stripped $RUST_SRC

ls -la binaries/
