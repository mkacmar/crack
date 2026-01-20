#!/bin/sh
set -ex

ARCH=$1
SRC=test/e2e/testdata/main.c
mkdir -p binaries

. test/e2e/testdata/log-env.sh

gcc -fPIE -pie -o binaries/${ARCH}-gcc-pie-explicit $SRC
gcc -fno-pie -no-pie -o binaries/${ARCH}-gcc-no-pie $SRC
gcc -static-pie -o binaries/${ARCH}-gcc-static-pie $SRC
gcc -shared -fPIC -o binaries/${ARCH}-gcc-shared $SRC
gcc -fPIE -pie -o binaries/${ARCH}-gcc-pie-stripped $SRC
strip binaries/${ARCH}-gcc-pie-stripped
# strip debug only, keep symbols
gcc -fPIE -pie -o binaries/${ARCH}-gcc-pie-strip-debug $SRC
strip --strip-debug binaries/${ARCH}-gcc-pie-strip-debug

clang -fPIE -pie -o binaries/${ARCH}-clang-pie-explicit $SRC
clang -fno-pie -no-pie -o binaries/${ARCH}-clang-no-pie $SRC
clang -fPIE -pie -o binaries/${ARCH}-clang-pie-stripped $SRC
strip binaries/${ARCH}-clang-pie-stripped

ls -la binaries/
