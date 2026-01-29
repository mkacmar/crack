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

gcc -D_FORTIFY_SOURCE=2 -O2 -o binaries/${ARCH}-gcc-fortify2-O2 $SRC
gcc -D_FORTIFY_SOURCE=1 -O1 -o binaries/${ARCH}-gcc-fortify1-O1 $SRC
gcc -D_FORTIFY_SOURCE=3 -O2 -o binaries/${ARCH}-gcc-fortify3-O2 $SRC || echo "FORTIFY_SOURCE=3 not supported"
gcc -U_FORTIFY_SOURCE -D_FORTIFY_SOURCE=0 -O2 -o binaries/${ARCH}-gcc-no-fortify $SRC
# -O0 disables fortify optimization
gcc -D_FORTIFY_SOURCE=2 -O0 -o binaries/${ARCH}-gcc-fortify2-O0 $SRC
gcc -D_FORTIFY_SOURCE=2 -O2 -o binaries/${ARCH}-gcc-fortify2-stripped $SRC
strip binaries/${ARCH}-gcc-fortify2-stripped
gcc -D_FORTIFY_SOURCE=2 -O2 -static -o binaries/${ARCH}-gcc-fortify2-static $SRC || echo "static linking not supported"
gcc -D_FORTIFY_SOURCE=2 -O2 -static -o binaries/${ARCH}-gcc-fortify2-static-stripped $SRC && \
  strip binaries/${ARCH}-gcc-fortify2-static-stripped || echo "static linking not supported"
gcc -D_FORTIFY_SOURCE=2 -O2 -flto -o binaries/${ARCH}-gcc-fortify2-lto $SRC
# simple program without fortifiable functions
gcc -D_FORTIFY_SOURCE=2 -O2 -o binaries/${ARCH}-gcc-fortify2-simple $SIMPLE

clang -D_FORTIFY_SOURCE=2 -O2 -o binaries/${ARCH}-clang-fortify2-O2 $SRC
clang -D_FORTIFY_SOURCE=1 -O1 -o binaries/${ARCH}-clang-fortify1-O1 $SRC
clang -U_FORTIFY_SOURCE -D_FORTIFY_SOURCE=0 -O2 -o binaries/${ARCH}-clang-no-fortify $SRC
# -O0 disables fortify optimization
clang -D_FORTIFY_SOURCE=2 -O0 -o binaries/${ARCH}-clang-fortify2-O0 $SRC
clang -D_FORTIFY_SOURCE=2 -O2 -o binaries/${ARCH}-clang-fortify2-stripped $SRC
strip binaries/${ARCH}-clang-fortify2-stripped
clang -D_FORTIFY_SOURCE=2 -O2 -flto -o binaries/${ARCH}-clang-fortify2-lto $SRC

ls -la binaries/
rm -f /tmp/fortify.c
