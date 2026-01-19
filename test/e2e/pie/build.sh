#!/bin/sh
set -ex

ARCH=$1
SRC=test/e2e/testdata/main.c
mkdir -p binaries

echo "=== Build environment ==="
uname -m
gcc --version | head -1
clang --version | head -1

gcc -fPIE -pie -o binaries/${ARCH}-gcc-pie-explicit $SRC
gcc -fno-pie -no-pie -o binaries/${ARCH}-gcc-no-pie $SRC
gcc -static-pie -o binaries/${ARCH}-gcc-static-pie $SRC
gcc -shared -fPIC -o binaries/${ARCH}-gcc-shared $SRC

clang -fPIE -pie -o binaries/${ARCH}-clang-pie-explicit $SRC
clang -fno-pie -no-pie -o binaries/${ARCH}-clang-no-pie $SRC

ls -la binaries/

