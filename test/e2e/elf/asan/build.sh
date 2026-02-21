#!/bin/sh
set -ex

ARCH=$1
C_SRC=test/e2e/elf/testdata/main.c
RUST_SRC=test/e2e/elf/testdata/main.rs
mkdir -p binaries

. test/e2e/elf/testdata/log-env.sh

build_c() { $1 $2 -o binaries/${ARCH}-$1-$3 $C_SRC; }
build_c_strip() { $1 $2 -o binaries/${ARCH}-$1-$3 $C_SRC && strip binaries/${ARCH}-$1-$3; }

build_c gcc "-fsanitize=address" asan
build_c_strip gcc "-fsanitize=address" asan-stripped

build_c clang "-fsanitize=address" asan
build_c_strip clang "-fsanitize=address" asan-stripped

build_c gcc "" no-asan
build_c clang "" no-asan

rustc -o binaries/${ARCH}-rustc-no-asan $RUST_SRC

if [ "$ARCH" != "arm" ]; then
    rustc +nightly -Zsanitizer=address -Cunsafe-allow-abi-mismatch=sanitizer -o binaries/${ARCH}-rustc-nightly-asan $RUST_SRC
    rustc +nightly -o binaries/${ARCH}-rustc-nightly-no-asan $RUST_SRC
fi

ls -la binaries/
