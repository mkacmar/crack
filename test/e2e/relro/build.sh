#!/bin/sh
set -ex

ARCH=$1
SRC=test/e2e/testdata/main.c
mkdir -p binaries

. test/e2e/testdata/log-env.sh

build() { $1 $2 -o binaries/${ARCH}-$1-$3 $SRC; }
build_strip() { $1 $2 -o binaries/${ARCH}-$1-$3 $SRC && strip binaries/${ARCH}-$1-$3; }

build gcc "-Wl,-z,relro" partial-relro
build gcc "-Wl,-z,relro,-z,now" full-relro
build gcc "-Wl,-z,norelro" no-relro
build_strip gcc "-Wl,-z,relro,-z,now" full-relro-stripped
build gcc "-static -Wl,-z,relro,-z,now" full-relro-static
build gcc "-shared -fPIC -Wl,-z,relro,-z,now" full-relro-shared
gcc -c -o binaries/${ARCH}-gcc-relocatable.o $SRC

build clang "-Wl,-z,relro" partial-relro
build clang "-Wl,-z,relro,-z,now" full-relro
build clang "-Wl,-z,norelro" no-relro
build_strip clang "-Wl,-z,relro,-z,now" full-relro-stripped
clang -c -o binaries/${ARCH}-clang-relocatable.o $SRC

ls -la binaries/
