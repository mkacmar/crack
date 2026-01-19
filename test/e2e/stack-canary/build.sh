#!/bin/sh
set -ex

ARCH=$1
mkdir -p binaries

echo "=== Build environment ==="
uname -m
gcc --version | head -1
clang --version | head -1

# Source with a buffer that triggers stack protection
cat > /tmp/vulnerable.c << 'EOF'
#include <string.h>
void vulnerable(const char *input) {
    char buffer[64];
    strcpy(buffer, input);
}
int main(void) {
    vulnerable("test");
    return 0;
}
EOF

# Simple source without vulnerable buffers (may not get canary with -fstack-protector)
cat > /tmp/simple.c << 'EOF'
int main(void) {
    return 0;
}
EOF

SRC=/tmp/vulnerable.c
SIMPLE=/tmp/simple.c

# --- GCC variants ---

# stack-protector-strong (recommended)
gcc -fstack-protector-strong -o binaries/${ARCH}-gcc-stack-protector-strong $SRC

# stack-protector-all (always adds canary)
gcc -fstack-protector-all -o binaries/${ARCH}-gcc-stack-protector-all $SRC

# stack-protector (basic, only some functions)
gcc -fstack-protector -o binaries/${ARCH}-gcc-stack-protector $SRC

# no stack protector
gcc -fno-stack-protector -o binaries/${ARCH}-gcc-no-stack-protector $SRC


# simple program with stack-protector (may not have canary - no vulnerable buffer)
gcc -fstack-protector -o binaries/${ARCH}-gcc-stack-protector-simple $SIMPLE

# simple program with stack-protector-all (forces canary even without buffer)
gcc -fstack-protector-all -o binaries/${ARCH}-gcc-stack-protector-all-simple $SIMPLE

# stripped binary with stack protection (tests dynamic symbol detection)
gcc -fstack-protector-strong -o binaries/${ARCH}-gcc-stack-protector-stripped $SRC
strip binaries/${ARCH}-gcc-stack-protector-stripped

# static binary with stack protection (tests static symbol detection)
gcc -fstack-protector-strong -static -o binaries/${ARCH}-gcc-stack-protector-static $SRC || echo "static linking not supported"

# static + stripped (edge case: both symbol tables may be empty - expected false negative)
gcc -fstack-protector-strong -static -o binaries/${ARCH}-gcc-stack-protector-static-stripped $SRC && \
  strip binaries/${ARCH}-gcc-stack-protector-static-stripped || echo "static linking not supported"

# LTO with stack protection (tests if canary survives link-time optimization)
gcc -flto -fstack-protector-strong -o binaries/${ARCH}-gcc-stack-protector-lto $SRC

# --- Clang variants ---

# stack-protector-strong
clang -fstack-protector-strong -o binaries/${ARCH}-clang-stack-protector-strong $SRC

# stack-protector-all
clang -fstack-protector-all -o binaries/${ARCH}-clang-stack-protector-all $SRC

# no stack protector
clang -fno-stack-protector -o binaries/${ARCH}-clang-no-stack-protector $SRC


# stripped binary with stack protection
clang -fstack-protector-strong -o binaries/${ARCH}-clang-stack-protector-stripped $SRC
strip binaries/${ARCH}-clang-stack-protector-stripped

# static binary with stack protection
clang -fstack-protector-strong -static -o binaries/${ARCH}-clang-stack-protector-static $SRC || echo "static linking not supported"

# static + stripped (edge case: both symbol tables may be empty - expected false negative)
clang -fstack-protector-strong -static -o binaries/${ARCH}-clang-stack-protector-static-stripped $SRC && \
  strip binaries/${ARCH}-clang-stack-protector-static-stripped || echo "static linking not supported"

# LTO with stack protection (tests if canary survives link-time optimization)
clang -flto -fstack-protector-strong -o binaries/${ARCH}-clang-stack-protector-lto $SRC

ls -la binaries/

# Cleanup
rm -f /tmp/vulnerable.c /tmp/simple.c

