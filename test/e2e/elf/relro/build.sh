#!/bin/sh
set -ex

ARCH=$1
C_SRC=test/e2e/elf/testdata/main.c
mkdir -p binaries

. test/e2e/elf/testdata/log-env.sh

build_c() { $1 $2 -o binaries/${ARCH}-$1-$3 $C_SRC; }
build_c_strip() { $1 $2 -o binaries/${ARCH}-$1-$3 $C_SRC && strip binaries/${ARCH}-$1-$3; }

build_c gcc "-Wl,-z,relro,-z,lazy" partial-relro
build_c gcc "-Wl,-z,relro,-z,now" full-relro
build_c gcc "-Wl,-z,norelro" no-relro
build_c_strip gcc "-Wl,-z,relro,-z,now" full-relro-stripped
build_c gcc "-static-pie -Wl,-z,relro,-z,now" full-relro-static
build_c gcc "-shared -fPIC -Wl,-z,relro,-z,now" full-relro-shared
gcc -c -o binaries/${ARCH}-gcc-relocatable.o $C_SRC

build_c clang "-Wl,-z,relro,-z,lazy" partial-relro
build_c clang "-Wl,-z,relro,-z,now" full-relro
build_c clang "-Wl,-z,norelro" no-relro
build_c_strip clang "-Wl,-z,relro,-z,now" full-relro-stripped
clang -c -o binaries/${ARCH}-clang-relocatable.o $C_SRC

ls -la binaries/
