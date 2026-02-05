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

build_c() { $1 $2 -o binaries/${ARCH}-$1-$3 $SRC; }
build_c_strip() { $1 $2 -o binaries/${ARCH}-$1-$3 $SRC && strip binaries/${ARCH}-$1-$3; }

CFI_FLAGS="-fsanitize=cfi -fsanitize-cfi-cross-dso -flto -fvisibility=hidden -fuse-ld=lld"

build_c clang "$CFI_FLAGS" cfi
build_c_strip clang "$CFI_FLAGS" cfi-stripped

build_c clang "" no-cfi
build_c_strip clang "" no-cfi-stripped

build_c gcc "" no-cfi
build_c_strip gcc "" no-cfi-stripped

ls -la binaries/
rm -f /tmp/cfi.c
