#!/bin/sh
set -ex

ARCH=$1
SRC=test/e2e/testdata/main.c
mkdir -p binaries

. test/e2e/testdata/log-env.sh

# partial RELRO (default on many systems)
gcc -Wl,-z,relro -o binaries/${ARCH}-gcc-partial-relro $SRC

# full RELRO
gcc -Wl,-z,relro,-z,now -o binaries/${ARCH}-gcc-full-relro $SRC

# no RELRO
gcc -Wl,-z,norelro -o binaries/${ARCH}-gcc-no-relro $SRC

# full RELRO stripped
gcc -Wl,-z,relro,-z,now -o binaries/${ARCH}-gcc-full-relro-stripped $SRC
strip binaries/${ARCH}-gcc-full-relro-stripped

# full RELRO static (static binaries can still have RELRO)
gcc -static -Wl,-z,relro,-z,now -o binaries/${ARCH}-gcc-full-relro-static $SRC

# shared library with full RELRO
gcc -shared -fPIC -Wl,-z,relro,-z,now -o binaries/${ARCH}-gcc-full-relro-shared $SRC

# partial RELRO
clang -Wl,-z,relro -o binaries/${ARCH}-clang-partial-relro $SRC

# full RELRO
clang -Wl,-z,relro,-z,now -o binaries/${ARCH}-clang-full-relro $SRC

# no RELRO
clang -Wl,-z,norelro -o binaries/${ARCH}-clang-no-relro $SRC

# full RELRO stripped
clang -Wl,-z,relro,-z,now -o binaries/${ARCH}-clang-full-relro-stripped $SRC
strip binaries/${ARCH}-clang-full-relro-stripped

ls -la binaries/

