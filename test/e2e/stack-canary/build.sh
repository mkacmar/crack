#!/bin/sh
set -ex

ARCH=$1
mkdir -p binaries

. test/e2e/testdata/log-env.sh

# source with a buffer that triggers stack protection
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

SRC=/tmp/vulnerable.c
SIMPLE=test/e2e/testdata/main.c

gcc -fstack-protector-strong -o binaries/${ARCH}-gcc-stack-protector-strong $SRC
gcc -fstack-protector-all -o binaries/${ARCH}-gcc-stack-protector-all $SRC
gcc -fstack-protector -o binaries/${ARCH}-gcc-stack-protector $SRC
gcc -fno-stack-protector -o binaries/${ARCH}-gcc-no-stack-protector $SRC
# simple program (no vulnerable buffer)
gcc -fstack-protector -o binaries/${ARCH}-gcc-stack-protector-simple $SIMPLE
gcc -fstack-protector-all -o binaries/${ARCH}-gcc-stack-protector-all-simple $SIMPLE
gcc -fstack-protector-strong -o binaries/${ARCH}-gcc-stack-protector-stripped $SRC
strip binaries/${ARCH}-gcc-stack-protector-stripped
gcc -fstack-protector-strong -static -o binaries/${ARCH}-gcc-stack-protector-static $SRC || echo "static linking not supported"
gcc -fstack-protector-strong -static -o binaries/${ARCH}-gcc-stack-protector-static-stripped $SRC && \
  strip binaries/${ARCH}-gcc-stack-protector-static-stripped || echo "static linking not supported"
gcc -flto -fstack-protector-strong -o binaries/${ARCH}-gcc-stack-protector-lto $SRC

clang -fstack-protector-strong -o binaries/${ARCH}-clang-stack-protector-strong $SRC
clang -fstack-protector-all -o binaries/${ARCH}-clang-stack-protector-all $SRC
clang -fno-stack-protector -o binaries/${ARCH}-clang-no-stack-protector $SRC
clang -fstack-protector-strong -o binaries/${ARCH}-clang-stack-protector-stripped $SRC
strip binaries/${ARCH}-clang-stack-protector-stripped
clang -fstack-protector-strong -static -o binaries/${ARCH}-clang-stack-protector-static $SRC || echo "static linking not supported"
clang -fstack-protector-strong -static -o binaries/${ARCH}-clang-stack-protector-static-stripped $SRC && \
  strip binaries/${ARCH}-clang-stack-protector-static-stripped || echo "static linking not supported"
clang -flto -fstack-protector-strong -o binaries/${ARCH}-clang-stack-protector-lto $SRC

ls -la binaries/
rm -f /tmp/vulnerable.c /tmp/simple.c
