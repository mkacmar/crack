#!/bin/sh
set -ex

ARCH=$1
mkdir -p binaries

. test/e2e/testdata/log-env.sh

cat > /tmp/ubsan.c << 'EOF'
int add(int a, int b) { return a + b; }
int divide(int a, int b) { return a / b; }

int main(void) {
    volatile int x = 100;
    return add(x, x) / divide(x, 1);
}
EOF

SRC=/tmp/ubsan.c

build() { $1 $2 -o binaries/${ARCH}-$1-$3 $SRC; }
build_strip() { $1 $2 -o binaries/${ARCH}-$1-$3 $SRC && strip binaries/${ARCH}-$1-$3; }

build gcc "-fsanitize=undefined" ubsan
build_strip gcc "-fsanitize=undefined" ubsan-stripped

build clang "-fsanitize=undefined" ubsan
build_strip clang "-fsanitize=undefined" ubsan-stripped

build gcc "" no-ubsan
build clang "" no-ubsan

ls -la binaries/
rm -f /tmp/ubsan.c
