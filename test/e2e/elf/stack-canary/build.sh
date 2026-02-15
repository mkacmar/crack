#!/bin/sh
set -ex

ARCH=$1
RUST_SRC=test/e2e/elf/testdata/main.rs
mkdir -p binaries

. test/e2e/elf/testdata/log-env.sh

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

C_SRC=/tmp/vulnerable.c
C_SRC_SIMPLE=test/e2e/elf/testdata/main.c

build_c() { $1 $2 -o binaries/${ARCH}-$1-$3 $4; }
build_c_strip() { $1 $2 -o binaries/${ARCH}-$1-$3 $4 && strip binaries/${ARCH}-$1-$3; }

build_c gcc -fstack-protector-strong stack-protector-strong $C_SRC
build_c gcc -fstack-protector-all stack-protector-all $C_SRC
build_c gcc -fstack-protector stack-protector $C_SRC
build_c gcc -fno-stack-protector no-stack-protector $C_SRC
build_c gcc -fstack-protector stack-protector-simple $C_SRC_SIMPLE
build_c gcc -fstack-protector-all stack-protector-all-simple $C_SRC_SIMPLE
build_c_strip gcc -fstack-protector-strong stack-protector-stripped $C_SRC
build_c gcc "-fstack-protector-strong -static" stack-protector-static $C_SRC
build_c_strip gcc "-fstack-protector-strong -static" stack-protector-static-stripped $C_SRC
build_c gcc "-flto -fstack-protector-strong" stack-protector-lto $C_SRC

build_c clang -fstack-protector-strong stack-protector-strong $C_SRC
build_c clang -fstack-protector-all stack-protector-all $C_SRC
build_c clang -fno-stack-protector no-stack-protector $C_SRC
build_c_strip clang -fstack-protector-strong stack-protector-stripped $C_SRC
build_c clang "-fstack-protector-strong -static" stack-protector-static $C_SRC
build_c_strip clang "-fstack-protector-strong -static" stack-protector-static-stripped $C_SRC
build_c clang "-flto -fstack-protector-strong" stack-protector-lto $C_SRC

rustc -o binaries/${ARCH}-rustc-stack-protector $RUST_SRC
rustc -C strip=symbols -o binaries/${ARCH}-rustc-stack-protector-stripped $RUST_SRC

ls -la binaries/
rm -f /tmp/vulnerable.c
