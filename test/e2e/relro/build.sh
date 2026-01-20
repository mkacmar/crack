#!/bin/sh
set -ex

ARCH=$1
SRC=test/e2e/testdata/main.c
mkdir -p binaries

. test/e2e/testdata/log-env.sh

gcc -Wl,-z,relro -o binaries/${ARCH}-gcc-partial-relro $SRC
gcc -Wl,-z,relro,-z,now -o binaries/${ARCH}-gcc-full-relro $SRC
gcc -Wl,-z,norelro -o binaries/${ARCH}-gcc-no-relro $SRC
gcc -Wl,-z,relro,-z,now -o binaries/${ARCH}-gcc-full-relro-stripped $SRC
strip binaries/${ARCH}-gcc-full-relro-stripped
gcc -static -Wl,-z,relro,-z,now -o binaries/${ARCH}-gcc-full-relro-static $SRC
gcc -shared -fPIC -Wl,-z,relro,-z,now -o binaries/${ARCH}-gcc-full-relro-shared $SRC

clang -Wl,-z,relro -o binaries/${ARCH}-clang-partial-relro $SRC
clang -Wl,-z,relro,-z,now -o binaries/${ARCH}-clang-full-relro $SRC
clang -Wl,-z,norelro -o binaries/${ARCH}-clang-no-relro $SRC
clang -Wl,-z,relro,-z,now -o binaries/${ARCH}-clang-full-relro-stripped $SRC
strip binaries/${ARCH}-clang-full-relro-stripped

ls -la binaries/
