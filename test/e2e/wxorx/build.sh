#!/bin/sh
set -ex

ARCH=$1
SRC=test/e2e/testdata/main.c
mkdir -p binaries

. test/e2e/testdata/log-env.sh

build() { $1 $2 -o binaries/${ARCH}-$1-$3 $SRC; }

build gcc "" wxorx
build gcc "-z execstack" execstack
build gcc "-shared -fPIC" shared-wxorx
build gcc "-shared -fPIC -z execstack" shared-execstack
gcc -c -o binaries/${ARCH}-gcc-relocatable.o $SRC

build clang "" wxorx
build clang "-z execstack" execstack
build clang "-shared -fPIC" shared-wxorx
build clang "-shared -fPIC -z execstack" shared-execstack
clang -c -o binaries/${ARCH}-clang-relocatable.o $SRC

ls -la binaries/
