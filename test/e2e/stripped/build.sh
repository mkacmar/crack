#!/bin/sh
set -ex

ARCH=$1
SRC=test/e2e/testdata/main.c
mkdir -p binaries

. test/e2e/testdata/log-env.sh

gcc -g -o binaries/${ARCH}-gcc-not-stripped $SRC

gcc -o binaries/${ARCH}-gcc-stripped $SRC
strip binaries/${ARCH}-gcc-stripped

gcc -g -o binaries/${ARCH}-gcc-strip-debug $SRC
strip --strip-debug binaries/${ARCH}-gcc-strip-debug

gcc -g -o binaries/${ARCH}-gcc-strip-symbols $SRC
strip --strip-unneeded binaries/${ARCH}-gcc-strip-symbols

gcc -s -o binaries/${ARCH}-gcc-link-stripped $SRC

gcc -g -o binaries/${ARCH}-gcc-partial-stripped $SRC
strip --strip-all --keep-section=.debug_info --keep-section=.debug_abbrev --keep-section=.debug_line binaries/${ARCH}-gcc-partial-stripped

clang -g -o binaries/${ARCH}-clang-not-stripped $SRC
clang -o binaries/${ARCH}-clang-stripped $SRC
strip binaries/${ARCH}-clang-stripped
clang -g -o binaries/${ARCH}-clang-strip-debug $SRC
strip --strip-debug binaries/${ARCH}-clang-strip-debug

ls -la binaries/
