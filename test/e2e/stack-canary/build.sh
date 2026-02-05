#!/bin/sh
set -ex

ARCH=$1
mkdir -p binaries

. test/e2e/testdata/log-env.sh

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

build_c() { $1 $2 -o binaries/${ARCH}-$1-$3 $4; }
build_c_strip() { $1 $2 -o binaries/${ARCH}-$1-$3 $4 && strip binaries/${ARCH}-$1-$3; }

build_c gcc -fstack-protector-strong stack-protector-strong $SRC
build_c gcc -fstack-protector-all stack-protector-all $SRC
build_c gcc -fstack-protector stack-protector $SRC
build_c gcc -fno-stack-protector no-stack-protector $SRC
build_c gcc -fstack-protector stack-protector-simple $SIMPLE
build_c gcc -fstack-protector-all stack-protector-all-simple $SIMPLE
build_c_strip gcc -fstack-protector-strong stack-protector-stripped $SRC
build_c gcc "-fstack-protector-strong -static" stack-protector-static $SRC
build_c_strip gcc "-fstack-protector-strong -static" stack-protector-static-stripped $SRC
build_c gcc "-flto -fstack-protector-strong" stack-protector-lto $SRC

build_c clang -fstack-protector-strong stack-protector-strong $SRC
build_c clang -fstack-protector-all stack-protector-all $SRC
build_c clang -fno-stack-protector no-stack-protector $SRC
build_c_strip clang -fstack-protector-strong stack-protector-stripped $SRC
build_c clang "-fstack-protector-strong -static" stack-protector-static $SRC
build_c_strip clang "-fstack-protector-strong -static" stack-protector-static-stripped $SRC
build_c clang "-flto -fstack-protector-strong" stack-protector-lto $SRC

ls -la binaries/
rm -f /tmp/vulnerable.c
