#!/bin/sh
set -ex

ARCH=$1
SRC=test/e2e/testdata/main.c
mkdir -p binaries

. test/e2e/testdata/log-env.sh

build_c() { $1 $2 -o binaries/${ARCH}-$1-$3 $SRC; }
build_c_strip() { $1 $2 -o binaries/${ARCH}-$1-$3 $SRC && strip binaries/${ARCH}-$1-$3; }

build_c gcc "-fPIE -pie -Wl,-z,noexecstack" aslr-full
build_c gcc "-fno-pie -no-pie" no-pie
build_c gcc "-fPIE -pie -Wl,-z,execstack" execstack
build_c gcc "-shared -fPIC" shared
build_c gcc "-static-pie -Wl,-z,noexecstack" static-pie
build_c_strip gcc "-fPIE -pie -Wl,-z,noexecstack" aslr-stripped
build_c gcc "-static -fno-pie -no-pie" static-no-pie

build_c clang "-fPIE -pie -Wl,-z,noexecstack" aslr-full
build_c clang "-fno-pie -no-pie" no-pie
build_c clang "-fPIE -pie -Wl,-z,execstack" execstack
build_c clang "-static -fno-pie -no-pie" static-no-pie

gcc -fPIE -pie -Wl,-z,noexecstack -o binaries/${ARCH}-gcc-textrel-patched $SRC
go run test/e2e/aslr/add-textrel.go binaries/${ARCH}-gcc-textrel-patched

gcc -c -o binaries/${ARCH}-gcc-object-file $SRC

ls -la binaries/
