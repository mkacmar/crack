#!/bin/sh
set -ex

ARCH=$1
SRC=test/e2e/testdata/main.c
mkdir -p binaries

. test/e2e/testdata/log-env.sh

build_c() { $1 $2 -o binaries/${ARCH}-$1-$3 $SRC; }

build_c gcc "" wxorx
build_c gcc "-z execstack" execstack
build_c gcc "-shared -fPIC" shared-wxorx
build_c gcc "-shared -fPIC -z execstack" shared-execstack
gcc -c -o binaries/${ARCH}-gcc-relocatable.o $SRC

build_c clang "" wxorx
build_c clang "-z execstack" execstack
build_c clang "-shared -fPIC" shared-wxorx
build_c clang "-shared -fPIC -z execstack" shared-execstack
clang -c -o binaries/${ARCH}-clang-relocatable.o $SRC

ls -la binaries/
