#!/bin/sh
set -ex

ARCH=$1
C_SRC=test/e2e/testdata/main.c
RUST_SRC=test/e2e/testdata/main.rs
mkdir -p binaries

. test/e2e/testdata/log-env.sh

build_c() { $1 -shared -fPIC $2 -o binaries/${ARCH}-$1-$3.so $C_SRC; }
build_c_strip() { build_c $1 "$2" $3 && strip binaries/${ARCH}-$1-$3.so; }

build_c gcc "-Wl,-z,nodlopen" nodlopen
build_c_strip gcc "-Wl,-z,nodlopen" nodlopen-stripped

build_c clang "-Wl,-z,nodlopen" nodlopen
build_c_strip clang "-Wl,-z,nodlopen" nodlopen-stripped

build_c gcc "" default
build_c clang "" default

gcc -fPIE -pie -o binaries/${ARCH}-gcc-pie-executable $C_SRC
clang -fPIE -pie -o binaries/${ARCH}-clang-pie-executable $C_SRC

rustc -o binaries/${ARCH}-rustc-executable $RUST_SRC
rustc --crate-type=cdylib -C link-arg=-z -C link-arg=nodlopen -o binaries/${ARCH}-rustc-nodlopen.so $RUST_SRC
rustc --crate-type=cdylib -C link-arg=-z -C link-arg=nodlopen -C strip=symbols -o binaries/${ARCH}-rustc-nodlopen-stripped.so $RUST_SRC
rustc --crate-type=cdylib -o binaries/${ARCH}-rustc-default.so $RUST_SRC

ls -la binaries/
