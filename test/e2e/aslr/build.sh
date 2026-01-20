#!/bin/sh
set -ex

ARCH=$1
SRC=test/e2e/testdata/main.c
mkdir -p binaries

. test/e2e/testdata/log-env.sh

# full ASLR compatible (PIE + NX stack)
gcc -fPIE -pie -Wl,-z,noexecstack -o binaries/${ARCH}-gcc-aslr-full $SRC
gcc -fno-pie -no-pie -o binaries/${ARCH}-gcc-no-pie $SRC
gcc -fPIE -pie -Wl,-z,execstack -o binaries/${ARCH}-gcc-execstack $SRC
gcc -shared -fPIC -o binaries/${ARCH}-gcc-shared $SRC
gcc -static-pie -Wl,-z,noexecstack -o binaries/${ARCH}-gcc-static-pie $SRC
gcc -fPIE -pie -Wl,-z,noexecstack -o binaries/${ARCH}-gcc-aslr-stripped $SRC
strip binaries/${ARCH}-gcc-aslr-stripped

clang -fPIE -pie -Wl,-z,noexecstack -o binaries/${ARCH}-clang-aslr-full $SRC
clang -fno-pie -no-pie -o binaries/${ARCH}-clang-no-pie $SRC
clang -fPIE -pie -Wl,-z,execstack -o binaries/${ARCH}-clang-execstack $SRC

ls -la binaries/


