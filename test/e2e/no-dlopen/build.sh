#!/bin/sh
set -ex

ARCH=$1
SRC=test/e2e/testdata/main.c
mkdir -p binaries

. test/e2e/testdata/log-env.sh

gcc -shared -fPIC -Wl,-z,nodlopen -o binaries/${ARCH}-gcc-nodlopen.so $SRC
gcc -shared -fPIC -Wl,-z,nodlopen -o binaries/${ARCH}-gcc-nodlopen-stripped.so $SRC
strip binaries/${ARCH}-gcc-nodlopen-stripped.so

clang -shared -fPIC -Wl,-z,nodlopen -o binaries/${ARCH}-clang-nodlopen.so $SRC
clang -shared -fPIC -Wl,-z,nodlopen -o binaries/${ARCH}-clang-nodlopen-stripped.so $SRC
strip binaries/${ARCH}-clang-nodlopen-stripped.so

gcc -shared -fPIC -o binaries/${ARCH}-gcc-default.so $SRC
clang -shared -fPIC -o binaries/${ARCH}-clang-default.so $SRC

ls -la binaries/
