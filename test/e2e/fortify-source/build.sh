#!/bin/sh
set -ex

ARCH=$1
mkdir -p binaries

. test/e2e/testdata/log-env.sh

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

SRC=/tmp/fortify.c
SIMPLE=test/e2e/testdata/main.c

build() { $1 $2 -o binaries/${ARCH}-$1-$3 $SRC; }
build_strip() { build "$@" && strip binaries/${ARCH}-$1-$3; }

build gcc "-D_FORTIFY_SOURCE=2 -O2" fortify2-O2
build gcc "-D_FORTIFY_SOURCE=1 -O1" fortify1-O1
build gcc "-D_FORTIFY_SOURCE=3 -O2" fortify3-O2
build gcc "-U_FORTIFY_SOURCE -D_FORTIFY_SOURCE=0 -O2" no-fortify
build gcc "-D_FORTIFY_SOURCE=2 -O0" fortify2-O0
build_strip gcc "-D_FORTIFY_SOURCE=2 -O2" fortify2-stripped
build gcc "-D_FORTIFY_SOURCE=2 -O2 -static" fortify2-static
build_strip gcc "-D_FORTIFY_SOURCE=2 -O2 -static" fortify2-static-stripped
build gcc "-D_FORTIFY_SOURCE=2 -O2 -flto" fortify2-lto
gcc -D_FORTIFY_SOURCE=2 -O2 -o binaries/${ARCH}-gcc-fortify2-simple $SIMPLE

build clang "-D_FORTIFY_SOURCE=2 -O2" fortify2-O2
build clang "-D_FORTIFY_SOURCE=1 -O1" fortify1-O1
build clang "-U_FORTIFY_SOURCE -D_FORTIFY_SOURCE=0 -O2" no-fortify
build clang "-D_FORTIFY_SOURCE=2 -O0" fortify2-O0
build_strip clang "-D_FORTIFY_SOURCE=2 -O2" fortify2-stripped
build clang "-D_FORTIFY_SOURCE=2 -O2 -flto" fortify2-lto

ls -la binaries/
rm -f /tmp/fortify.c
