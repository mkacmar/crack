#!/bin/sh
set -ex

ARCH=$1
C_SRC=test/e2e/testdata/main.c
RUST_SRC=test/e2e/testdata/main.rs
mkdir -p binaries

. test/e2e/testdata/log-env.sh

build_c() { $1 $2 -o binaries/${ARCH}-$1-$3 $C_SRC; }
build_c_strip() { $1 $2 -o binaries/${ARCH}-$1-$3 $C_SRC && strip binaries/${ARCH}-$1-$3; }

build_c gcc "-fPIE -pie" pie-explicit
build_c gcc "-fno-pie -no-pie" no-pie
build_c gcc -static-pie static-pie
build_c gcc "-shared -fPIC" shared
build_c_strip gcc "-fPIE -pie" pie-stripped
gcc -fPIE -pie -o binaries/${ARCH}-gcc-pie-strip-debug $C_SRC && strip --strip-debug binaries/${ARCH}-gcc-pie-strip-debug

build_c clang "-fPIE -pie" pie-explicit
build_c clang "-fno-pie -no-pie" no-pie
build_c_strip clang "-fPIE -pie" pie-stripped

gcc -c -o binaries/${ARCH}-gcc-relocatable.o $C_SRC

rustc -C relocation-model=pic -o binaries/${ARCH}-rustc-pie $RUST_SRC
rustc -C relocation-model=static -o binaries/${ARCH}-rustc-no-pie $RUST_SRC
rustc -C relocation-model=pic -C strip=symbols -o binaries/${ARCH}-rustc-pie-stripped $RUST_SRC

ls -la binaries/
