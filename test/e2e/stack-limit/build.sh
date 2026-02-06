#!/bin/sh
set -ex

ARCH=$1
C_SRC=test/e2e/testdata/main.c
RUST_SRC=test/e2e/testdata/main.rs
mkdir -p binaries

. test/e2e/testdata/log-env.sh

build_c() { $1 $2 -o binaries/${ARCH}-$1-$3 $C_SRC; }

build_c gcc "-Wl,-z,stack-size=8388608" stack-limit
build_c gcc "" no-stack-limit
gcc -c -o binaries/${ARCH}-gcc-relocatable.o $C_SRC

build_c clang "-Wl,-z,stack-size=8388608" stack-limit
build_c clang "" no-stack-limit
clang -c -o binaries/${ARCH}-clang-relocatable.o $C_SRC

rustc -o binaries/${ARCH}-rustc-no-stack-limit $RUST_SRC

ls -la binaries/
