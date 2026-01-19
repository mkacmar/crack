#!/bin/sh
set -ex

ARCH=$1
SRC=test/e2e/testdata/main.c
mkdir -p binaries

echo "=== Build environment ==="
uname -m
gcc --version | head -1
clang --version | head -1

# NX explicitly enabled
gcc -Wl,-z,noexecstack -o binaries/${ARCH}-gcc-nx-explicit $SRC

# NX disabled (executable stack)
gcc -Wl,-z,execstack -o binaries/${ARCH}-gcc-no-nx $SRC

# stripped with NX enabled
gcc -Wl,-z,noexecstack -o binaries/${ARCH}-gcc-nx-stripped $SRC
strip binaries/${ARCH}-gcc-nx-stripped

# static with NX enabled
gcc -Wl,-z,noexecstack -static -o binaries/${ARCH}-gcc-nx-static $SRC || echo "static linking not supported"

# NX explicitly enabled
clang -Wl,-z,noexecstack -o binaries/${ARCH}-clang-nx-explicit $SRC

# NX disabled (executable stack)
clang -Wl,-z,execstack -o binaries/${ARCH}-clang-no-nx $SRC

# stripped with NX enabled
clang -Wl,-z,noexecstack -o binaries/${ARCH}-clang-nx-stripped $SRC
strip binaries/${ARCH}-clang-nx-stripped

ls -la binaries/

