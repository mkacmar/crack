#!/bin/sh
set -ex

ARCH=$1
SRC=test/e2e/testdata/main.c
mkdir -p binaries

. test/e2e/testdata/log-env.sh

build() { $1 $2 -o binaries/${ARCH}-$1-$3 $SRC; }
build_strip() { $1 $2 -o binaries/${ARCH}-$1-$3 $SRC && strip binaries/${ARCH}-$1-$3; }

build gcc "-fPIE -pie -Wl,-z,noexecstack" aslr-full
build gcc "-fno-pie -no-pie" no-pie
build gcc "-fPIE -pie -Wl,-z,execstack" execstack
build gcc "-shared -fPIC" shared
build gcc "-static-pie -Wl,-z,noexecstack" static-pie
build_strip gcc "-fPIE -pie -Wl,-z,noexecstack" aslr-stripped
build gcc "-static -fno-pie -no-pie" static-no-pie

build clang "-fPIE -pie -Wl,-z,noexecstack" aslr-full
build clang "-fno-pie -no-pie" no-pie
build clang "-fPIE -pie -Wl,-z,execstack" execstack
build clang "-static -fno-pie -no-pie" static-no-pie

gcc -fPIE -pie -Wl,-z,noexecstack -o binaries/${ARCH}-gcc-textrel-patched $SRC
go run test/e2e/aslr/add-textrel.go binaries/${ARCH}-gcc-textrel-patched

gcc -c -o binaries/${ARCH}-gcc-object-file $SRC

ls -la binaries/
