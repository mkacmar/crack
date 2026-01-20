#!/bin/sh
set -ex

ARCH=$1
SRC=test/e2e/testdata/main.c
mkdir -p binaries

. test/e2e/testdata/log-env.sh

# not stripped (has symbols and debug info)
gcc -g -o binaries/${ARCH}-gcc-not-stripped $SRC
# fully stripped
gcc -o binaries/${ARCH}-gcc-stripped $SRC
strip binaries/${ARCH}-gcc-stripped
# strip debug only (keeps symbol table)
gcc -g -o binaries/${ARCH}-gcc-strip-debug $SRC
strip --strip-debug binaries/${ARCH}-gcc-strip-debug
# strip symbols only (keeps debug info) - unusual but possible
gcc -g -o binaries/${ARCH}-gcc-strip-symbols $SRC
strip --strip-unneeded binaries/${ARCH}-gcc-strip-symbols
# compiled with -s (stripped at link time)
gcc -s -o binaries/${ARCH}-gcc-link-stripped $SRC

clang -g -o binaries/${ARCH}-clang-not-stripped $SRC
clang -o binaries/${ARCH}-clang-stripped $SRC
strip binaries/${ARCH}-clang-stripped
clang -g -o binaries/${ARCH}-clang-strip-debug $SRC
strip --strip-debug binaries/${ARCH}-clang-strip-debug

ls -la binaries/

