#!/bin/sh
set -ex

ARCH=$1
SRC=test/e2e/testdata/main.c
mkdir -p binaries

. test/e2e/testdata/log-env.sh

gcc -fPIE -pie -Wl,-z,noexecstack -o binaries/${ARCH}-gcc-aslr-full $SRC
gcc -fno-pie -no-pie -o binaries/${ARCH}-gcc-no-pie $SRC
gcc -fPIE -pie -Wl,-z,execstack -o binaries/${ARCH}-gcc-execstack $SRC
gcc -shared -fPIC -o binaries/${ARCH}-gcc-shared $SRC
gcc -static-pie -Wl,-z,noexecstack -o binaries/${ARCH}-gcc-static-pie $SRC
gcc -fPIE -pie -Wl,-z,noexecstack -o binaries/${ARCH}-gcc-aslr-stripped $SRC
strip binaries/${ARCH}-gcc-aslr-stripped
gcc -static -fno-pie -no-pie -o binaries/${ARCH}-gcc-static-no-pie $SRC

clang -fPIE -pie -Wl,-z,noexecstack -o binaries/${ARCH}-clang-aslr-full $SRC
clang -fno-pie -no-pie -o binaries/${ARCH}-clang-no-pie $SRC
clang -fPIE -pie -Wl,-z,execstack -o binaries/${ARCH}-clang-execstack $SRC
clang -static -fno-pie -no-pie -o binaries/${ARCH}-clang-static-no-pie $SRC

# Text relocations binary (patch DT_DEBUG -> DT_TEXTREL)
gcc -fPIE -pie -Wl,-z,noexecstack -o binaries/${ARCH}-gcc-textrel-patched $SRC
go run test/e2e/aslr/add-textrel.go binaries/${ARCH}-gcc-textrel-patched

ls -la binaries/


