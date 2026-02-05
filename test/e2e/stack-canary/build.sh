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

build() { $1 $2 -o binaries/${ARCH}-$1-$3 $4; }
build_strip() { $1 $2 -o binaries/${ARCH}-$1-$3 $4 && strip binaries/${ARCH}-$1-$3; }

build gcc -fstack-protector-strong stack-protector-strong $SRC
build gcc -fstack-protector-all stack-protector-all $SRC
build gcc -fstack-protector stack-protector $SRC
build gcc -fno-stack-protector no-stack-protector $SRC
build gcc -fstack-protector stack-protector-simple $SIMPLE
build gcc -fstack-protector-all stack-protector-all-simple $SIMPLE
build_strip gcc -fstack-protector-strong stack-protector-stripped $SRC
build gcc "-fstack-protector-strong -static" stack-protector-static $SRC
build_strip gcc "-fstack-protector-strong -static" stack-protector-static-stripped $SRC
build gcc "-flto -fstack-protector-strong" stack-protector-lto $SRC

build clang -fstack-protector-strong stack-protector-strong $SRC
build clang -fstack-protector-all stack-protector-all $SRC
build clang -fno-stack-protector no-stack-protector $SRC
build_strip clang -fstack-protector-strong stack-protector-stripped $SRC
build clang "-fstack-protector-strong -static" stack-protector-static $SRC
build_strip clang "-fstack-protector-strong -static" stack-protector-static-stripped $SRC
build clang "-flto -fstack-protector-strong" stack-protector-lto $SRC

ls -la binaries/
rm -f /tmp/vulnerable.c
