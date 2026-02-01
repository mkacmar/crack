#!/bin/sh
set -ex

ARCH=$1
SRC=test/e2e/testdata/main.c
mkdir -p binaries

. test/e2e/testdata/log-env.sh

gcc -o binaries/${ARCH}-gcc-wxorx $SRC
gcc -z execstack -o binaries/${ARCH}-gcc-execstack $SRC
gcc -shared -fPIC -o binaries/${ARCH}-gcc-shared-wxorx $SRC
gcc -shared -fPIC -z execstack -o binaries/${ARCH}-gcc-shared-execstack $SRC
gcc -c -o binaries/${ARCH}-gcc-relocatable.o $SRC

clang -o binaries/${ARCH}-clang-wxorx $SRC
clang -z execstack -o binaries/${ARCH}-clang-execstack $SRC
clang -shared -fPIC -o binaries/${ARCH}-clang-shared-wxorx $SRC
clang -shared -fPIC -z execstack -o binaries/${ARCH}-clang-shared-execstack $SRC
clang -c -o binaries/${ARCH}-clang-relocatable.o $SRC

ls -la binaries/
