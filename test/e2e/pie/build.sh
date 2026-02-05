#!/bin/sh
set -ex

ARCH=$1
SRC=test/e2e/testdata/main.c
mkdir -p binaries

. test/e2e/testdata/log-env.sh

build() { $1 $2 -o binaries/${ARCH}-$1-$3 $SRC; }
build_strip() { $1 $2 -o binaries/${ARCH}-$1-$3 $SRC && strip binaries/${ARCH}-$1-$3; }

build gcc "-fPIE -pie" pie-explicit
build gcc "-fno-pie -no-pie" no-pie
build gcc -static-pie static-pie
build gcc "-shared -fPIC" shared
build_strip gcc "-fPIE -pie" pie-stripped
gcc -fPIE -pie -o binaries/${ARCH}-gcc-pie-strip-debug $SRC && strip --strip-debug binaries/${ARCH}-gcc-pie-strip-debug

build clang "-fPIE -pie" pie-explicit
build clang "-fno-pie -no-pie" no-pie
build_strip clang "-fPIE -pie" pie-stripped

gcc -c -o binaries/${ARCH}-gcc-object-file $SRC

ls -la binaries/
