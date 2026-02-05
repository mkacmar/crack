#!/bin/sh
set -ex

ARCH=$1
SRC=test/e2e/testdata/main.c
mkdir -p binaries

. test/e2e/testdata/log-env.sh

build() { $1 -shared -fPIC $2 -o binaries/${ARCH}-$1-$3.so $SRC; }
build_strip() { build $1 "$2" $3 && strip binaries/${ARCH}-$1-$3.so; }

build gcc "-Wl,-z,nodlopen" nodlopen
build_strip gcc "-Wl,-z,nodlopen" nodlopen-stripped

build clang "-Wl,-z,nodlopen" nodlopen
build_strip clang "-Wl,-z,nodlopen" nodlopen-stripped

build gcc "" default
build clang "" default

gcc -fPIE -pie -o binaries/${ARCH}-gcc-pie-executable $SRC
clang -fPIE -pie -o binaries/${ARCH}-clang-pie-executable $SRC

ls -la binaries/
