#!/bin/sh
set -ex

ARCH=$1
C_SRC=test/e2e/testdata/main.c
RUST_SRC=test/e2e/testdata/main.rs
mkdir -p binaries

. test/e2e/testdata/log-env.sh

build_c() { $1 $2 -o binaries/${ARCH}-$1-$3 $C_SRC; }
build_c_strip() { $1 $2 -o binaries/${ARCH}-$1-$3 $C_SRC && strip binaries/${ARCH}-$1-$3; }

build_c clang "-fsanitize=safe-stack" safestack
build_c_strip clang "-fsanitize=safe-stack" safestack-stripped

build_c clang "" no-safestack
build_c gcc "" no-safestack

rustc -o binaries/${ARCH}-rustc-no-safestack $RUST_SRC

ls -la binaries/
