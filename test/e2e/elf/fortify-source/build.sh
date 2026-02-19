#!/bin/sh
set -ex

ARCH=$1
RUST_SRC=test/e2e/elf/testdata/main.rs
mkdir -p binaries

. test/e2e/elf/testdata/log-env.sh

cat > /tmp/fortify.c << 'EOF'
#include <string.h>
#include <stdio.h>

int main(int argc, char **argv) {
    char buf[64];
    char *src = argv[0];
    size_t len = strlen(src);

    if (len > 32) len = 32;
    memcpy(buf, src, len);
    strcpy(buf, src);

    puts(buf);
    return 0;
}
EOF

C_SRC=/tmp/fortify.c
C_SRC_SIMPLE=test/e2e/elf/testdata/main.c

build_c() { $1 $2 -o binaries/${ARCH}-$1-$3 $C_SRC; }
build_c_strip() { build_c "$@" && strip binaries/${ARCH}-$1-$3; }

build_c gcc "-D_FORTIFY_SOURCE=2 -O2" fortify2-O2
build_c gcc "-D_FORTIFY_SOURCE=1 -O1" fortify1-O1
build_c gcc "-D_FORTIFY_SOURCE=3 -O2" fortify3-O2
build_c gcc "-U_FORTIFY_SOURCE -D_FORTIFY_SOURCE=0 -O2" no-fortify
build_c gcc "-D_FORTIFY_SOURCE=2 -O0" fortify2-O0
build_c_strip gcc "-D_FORTIFY_SOURCE=2 -O2" fortify2-stripped
build_c gcc "-D_FORTIFY_SOURCE=2 -O2 -static" fortify2-static
build_c_strip gcc "-D_FORTIFY_SOURCE=2 -O2 -static" fortify2-static-stripped
build_c gcc "-D_FORTIFY_SOURCE=2 -O2 -flto" fortify2-lto
gcc -D_FORTIFY_SOURCE=2 -O2 -o binaries/${ARCH}-gcc-fortify2-simple $C_SRC_SIMPLE

build_c clang "-D_FORTIFY_SOURCE=2 -O2" fortify2-O2
build_c clang "-D_FORTIFY_SOURCE=1 -O1" fortify1-O1
build_c clang "-U_FORTIFY_SOURCE -D_FORTIFY_SOURCE=0 -O2" no-fortify
build_c clang "-D_FORTIFY_SOURCE=2 -O0" fortify2-O0
build_c_strip clang "-D_FORTIFY_SOURCE=2 -O2" fortify2-stripped
build_c clang "-D_FORTIFY_SOURCE=2 -O2 -flto" fortify2-lto

build_c gcc "-shared -fPIC -D_FORTIFY_SOURCE=2 -O2" fortify2-shared
build_c gcc "-shared -fPIC -U_FORTIFY_SOURCE -D_FORTIFY_SOURCE=0 -O2" no-fortify-shared

rustc -o binaries/${ARCH}-rustc-no-fortify $RUST_SRC

ls -la binaries/
rm -f /tmp/fortify.c
