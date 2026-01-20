#!/bin/sh
set -ex

ARCH=$1
mkdir -p binaries

echo "=== Build environment ==="
uname -m
gcc --version | head -1
clang --version | head -1

# Source that uses fortifiable functions with runtime-determined sizes
# Using strlen() forces runtime evaluation, preventing compile-time optimization
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

# --- GCC variants ---

# FORTIFY_SOURCE=2 with -O2 (common production setting)
gcc -D_FORTIFY_SOURCE=2 -O2 -o binaries/${ARCH}-gcc-fortify2-O2 $SRC

# FORTIFY_SOURCE=1 with -O1
gcc -D_FORTIFY_SOURCE=1 -O1 -o binaries/${ARCH}-gcc-fortify1-O1 $SRC

# FORTIFY_SOURCE=3 with -O2 (strongest, GCC 12+)
gcc -D_FORTIFY_SOURCE=3 -O2 -o binaries/${ARCH}-gcc-fortify3-O2 $SRC || echo "FORTIFY_SOURCE=3 not supported"

# no FORTIFY_SOURCE with -O2
gcc -O2 -o binaries/${ARCH}-gcc-no-fortify $SRC

# FORTIFY_SOURCE=2 but -O0 (fortify ignored without optimization)
gcc -D_FORTIFY_SOURCE=2 -O0 -o binaries/${ARCH}-gcc-fortify2-O0 $SRC

# stripped with FORTIFY_SOURCE
gcc -D_FORTIFY_SOURCE=2 -O2 -o binaries/${ARCH}-gcc-fortify2-stripped $SRC
strip binaries/${ARCH}-gcc-fortify2-stripped

# static with FORTIFY_SOURCE
gcc -D_FORTIFY_SOURCE=2 -O2 -static -o binaries/${ARCH}-gcc-fortify2-static $SRC || echo "static linking not supported"

# static + stripped (edge case: false negative expected)
gcc -D_FORTIFY_SOURCE=2 -O2 -static -o binaries/${ARCH}-gcc-fortify2-static-stripped $SRC && \
  strip binaries/${ARCH}-gcc-fortify2-static-stripped || echo "static linking not supported"

# LTO with FORTIFY_SOURCE
gcc -D_FORTIFY_SOURCE=2 -O2 -flto -o binaries/${ARCH}-gcc-fortify2-lto $SRC

# simple program without fortifiable functions (should skip)
gcc -D_FORTIFY_SOURCE=2 -O2 -o binaries/${ARCH}-gcc-fortify2-simple $SIMPLE

# --- Clang variants ---

# FORTIFY_SOURCE=2 with -O2
clang -D_FORTIFY_SOURCE=2 -O2 -o binaries/${ARCH}-clang-fortify2-O2 $SRC

# FORTIFY_SOURCE=1 with -O1
clang -D_FORTIFY_SOURCE=1 -O1 -o binaries/${ARCH}-clang-fortify1-O1 $SRC

# no FORTIFY_SOURCE with -O2
clang -O2 -o binaries/${ARCH}-clang-no-fortify $SRC

# FORTIFY_SOURCE=2 but -O0 (fortify ignored without optimization)
clang -D_FORTIFY_SOURCE=2 -O0 -o binaries/${ARCH}-clang-fortify2-O0 $SRC

# stripped with FORTIFY_SOURCE
clang -D_FORTIFY_SOURCE=2 -O2 -o binaries/${ARCH}-clang-fortify2-stripped $SRC
strip binaries/${ARCH}-clang-fortify2-stripped

# LTO with FORTIFY_SOURCE
clang -D_FORTIFY_SOURCE=2 -O2 -flto -o binaries/${ARCH}-clang-fortify2-lto $SRC

ls -la binaries/

# Cleanup
rm -f /tmp/fortify.c

