#!/bin/sh
set -ex

ARCH=$1
SRC=test/e2e/testdata/main.c
mkdir -p binaries

echo "=== Build environment ==="
uname -m
gcc --version | head -1
clang --version | head -1

# NX enabled (default)
gcc -o binaries/${ARCH}-gcc-nx-default $SRC

# NX explicitly enabled
gcc -Wl,-z,noexecstack -o binaries/${ARCH}-gcc-nx-explicit $SRC

# NX disabled (executable stack)
gcc -Wl,-z,execstack -o binaries/${ARCH}-gcc-no-nx $SRC

# NX enabled (default)
clang -o binaries/${ARCH}-clang-nx-default $SRC

# NX explicitly enabled
clang -Wl,-z,noexecstack -o binaries/${ARCH}-clang-nx-explicit $SRC

# NX disabled (executable stack)
clang -Wl,-z,execstack -o binaries/${ARCH}-clang-no-nx $SRC

ls -la binaries/

