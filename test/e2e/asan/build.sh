#!/bin/sh
set -ex

ARCH=$1
SRC=test/e2e/testdata/main.c
mkdir -p binaries

. test/e2e/testdata/log-env.sh

build() { $1 $2 -o binaries/${ARCH}-$1-$3 $SRC; }
build_strip() { $1 $2 -o binaries/${ARCH}-$1-$3 $SRC && strip binaries/${ARCH}-$1-$3; }

build gcc "-fsanitize=address" asan
build_strip gcc "-fsanitize=address" asan-stripped

build clang "-fsanitize=address" asan
build_strip clang "-fsanitize=address" asan-stripped

build gcc "" no-asan
build clang "" no-asan

ls -la binaries/
