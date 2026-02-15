#!/bin/sh
set -ex

ARCH=$1
C_SRC=test/e2e/elf/testdata/main.c
RUST_SRC=test/e2e/elf/testdata/main.rs
mkdir -p binaries

. test/e2e/elf/testdata/log-env.sh

gcc -g -o binaries/${ARCH}-gcc-not-stripped $C_SRC

gcc -o binaries/${ARCH}-gcc-stripped $C_SRC
strip binaries/${ARCH}-gcc-stripped

gcc -g -o binaries/${ARCH}-gcc-strip-debug $C_SRC
strip --strip-debug binaries/${ARCH}-gcc-strip-debug

gcc -g -o binaries/${ARCH}-gcc-strip-symbols $C_SRC
strip --strip-unneeded binaries/${ARCH}-gcc-strip-symbols

gcc -s -o binaries/${ARCH}-gcc-link-stripped $C_SRC

gcc -g -o binaries/${ARCH}-gcc-partial-stripped $C_SRC
strip --strip-all --keep-section=.debug_info --keep-section=.debug_abbrev --keep-section=.debug_line binaries/${ARCH}-gcc-partial-stripped

clang -g -o binaries/${ARCH}-clang-not-stripped $C_SRC
clang -o binaries/${ARCH}-clang-stripped $C_SRC
strip binaries/${ARCH}-clang-stripped
clang -g -o binaries/${ARCH}-clang-strip-debug $C_SRC
strip --strip-debug binaries/${ARCH}-clang-strip-debug

rustc -o binaries/${ARCH}-rustc-not-stripped $RUST_SRC
rustc -C strip=symbols -o binaries/${ARCH}-rustc-stripped $RUST_SRC
rustc -C strip=debuginfo -o binaries/${ARCH}-rustc-strip-debuginfo $RUST_SRC

ls -la binaries/
