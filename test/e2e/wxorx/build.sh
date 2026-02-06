#!/bin/sh
set -ex

ARCH=$1
C_SRC=test/e2e/testdata/main.c
RUST_SRC=test/e2e/testdata/main.rs
mkdir -p binaries

. test/e2e/testdata/log-env.sh

build_c() { $1 $2 -o binaries/${ARCH}-$1-$3 $C_SRC; }

build_c gcc "" wxorx
build_c gcc "-z execstack" execstack
build_c gcc "-shared -fPIC" shared-wxorx
build_c gcc "-shared -fPIC -z execstack" shared-execstack
gcc -c -o binaries/${ARCH}-gcc-relocatable.o $C_SRC

build_c clang "" wxorx
build_c clang "-z execstack" execstack
build_c clang "-shared -fPIC" shared-wxorx
build_c clang "-shared -fPIC -z execstack" shared-execstack
clang -c -o binaries/${ARCH}-clang-relocatable.o $C_SRC

rustc -o binaries/${ARCH}-rustc-wxorx $RUST_SRC

ls -la binaries/
