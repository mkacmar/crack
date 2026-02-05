#!/bin/sh
set -ex

ARCH=$1
SRC=test/e2e/testdata/main.c
mkdir -p binaries

. test/e2e/testdata/log-env.sh

build_c() { $1 $2 -o binaries/${ARCH}-$1-$3 $SRC; }
build_c_strip() { $1 $2 -o binaries/${ARCH}-$1-$3 $SRC && strip binaries/${ARCH}-$1-$3; }

build_c clang "-fsanitize=safe-stack" safestack
build_c_strip clang "-fsanitize=safe-stack" safestack-stripped

build_c clang "" no-safestack
build_c gcc "" no-safestack

ls -la binaries/
