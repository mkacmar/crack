#!/bin/sh
set -ex

ARCH=$1
RUST_SRC=test/e2e/elf/testdata/main.rs
mkdir -p binaries

. test/e2e/elf/testdata/log-env.sh

cat > /tmp/ubsan.c << 'EOF'
int add(int a, int b) { return a + b; }
int divide(int a, int b) { return a / b; }

int main(void) {
    volatile int x = 100;
    return add(x, x) / divide(x, 1);
}
EOF

C_SRC=/tmp/ubsan.c

build_c() { $1 $2 -o binaries/${ARCH}-$1-$3 $C_SRC; }
build_c_strip() { $1 $2 -o binaries/${ARCH}-$1-$3 $C_SRC && strip binaries/${ARCH}-$1-$3; }

build_c gcc "-fsanitize=undefined" ubsan
build_c_strip gcc "-fsanitize=undefined" ubsan-stripped

build_c clang "-fsanitize=undefined" ubsan
build_c_strip clang "-fsanitize=undefined" ubsan-stripped

build_c gcc "" no-ubsan
build_c clang "" no-ubsan

rustc -o binaries/${ARCH}-rustc-no-ubsan $RUST_SRC

ls -la binaries/
rm -f /tmp/ubsan.c
