#!/bin/sh
set -ex

ARCH=$1
SRC=test/e2e/testdata/main.c
mkdir -p binaries

. test/e2e/testdata/log-env.sh

build_c() { $1 -shared -fPIC $2 -o binaries/${ARCH}-$1-$3.so $SRC; }
build_c_strip() { build $1 "$2" $3 && strip binaries/${ARCH}-$1-$3.so; }

build_c gcc "-Wl,-z,nodlopen" nodlopen
build_c_strip gcc "-Wl,-z,nodlopen" nodlopen-stripped

build_c clang "-Wl,-z,nodlopen" nodlopen
build_c_strip clang "-Wl,-z,nodlopen" nodlopen-stripped

build_c gcc "" default
build_c clang "" default

gcc -fPIE -pie -o binaries/${ARCH}-gcc-pie-executable $SRC
clang -fPIE -pie -o binaries/${ARCH}-clang-pie-executable $SRC

ls -la binaries/
