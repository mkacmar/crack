#!/bin/sh
set -ex

ARCH=$1
mkdir -p binaries

. test/e2e/testdata/log-env.sh

cat > /tmp/cfi.c << 'EOF'
typedef int (*func_ptr)(int);

int add_one(int x) { return x + 1; }

int main(void) {
    func_ptr fn = add_one;
    return fn(5) > 0 ? 0 : 1;
}
EOF

SRC=/tmp/cfi.c

# Cross-DSO CFI with LTO
clang -fsanitize=cfi -fsanitize-cfi-cross-dso -flto -fvisibility=hidden -fuse-ld=lld -o binaries/${ARCH}-clang-cfi $SRC
clang -fsanitize=cfi -fsanitize-cfi-cross-dso -flto -fvisibility=hidden -fuse-ld=lld -o binaries/${ARCH}-clang-cfi-stripped $SRC
strip binaries/${ARCH}-clang-cfi-stripped

# No CFI
clang -o binaries/${ARCH}-clang-no-cfi $SRC
clang -o binaries/${ARCH}-clang-no-cfi-stripped $SRC
strip binaries/${ARCH}-clang-no-cfi-stripped

gcc -o binaries/${ARCH}-gcc-no-cfi $SRC
gcc -o binaries/${ARCH}-gcc-no-cfi-stripped $SRC
strip binaries/${ARCH}-gcc-no-cfi-stripped

ls -la binaries/
rm -f /tmp/cfi.c
