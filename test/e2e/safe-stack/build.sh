#!/bin/sh
set -ex

ARCH=$1
SRC=test/e2e/testdata/main.c
mkdir -p binaries

. test/e2e/testdata/log-env.sh

build() { $1 $2 -o binaries/${ARCH}-$1-$3 $SRC; }
build_strip() { $1 $2 -o binaries/${ARCH}-$1-$3 $SRC && strip binaries/${ARCH}-$1-$3; }

build clang "-fsanitize=safe-stack" safestack
build_strip clang "-fsanitize=safe-stack" safestack-stripped

build clang "" no-safestack
build gcc "" no-safestack

ls -la binaries/
