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

gcc -fsanitize=undefined -o binaries/${ARCH}-gcc-ubsan $SRC
gcc -fsanitize=undefined -o binaries/${ARCH}-gcc-ubsan-stripped $SRC
strip binaries/${ARCH}-gcc-ubsan-stripped

clang -fsanitize=undefined -o binaries/${ARCH}-clang-ubsan $SRC
clang -fsanitize=undefined -o binaries/${ARCH}-clang-ubsan-stripped $SRC
strip binaries/${ARCH}-clang-ubsan-stripped

gcc -o binaries/${ARCH}-gcc-no-ubsan $SRC
clang -o binaries/${ARCH}-clang-no-ubsan $SRC

ls -la binaries/
rm -f /tmp/ubsan.c
