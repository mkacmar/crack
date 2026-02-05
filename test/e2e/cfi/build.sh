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

build() { $1 $2 -o binaries/${ARCH}-$1-$3 $SRC; }
build_strip() { $1 $2 -o binaries/${ARCH}-$1-$3 $SRC && strip binaries/${ARCH}-$1-$3; }

CFI_FLAGS="-fsanitize=cfi -fsanitize-cfi-cross-dso -flto -fvisibility=hidden -fuse-ld=lld"

build clang "$CFI_FLAGS" cfi
build_strip clang "$CFI_FLAGS" cfi-stripped

build clang "" no-cfi
build_strip clang "" no-cfi-stripped

build gcc "" no-cfi
build_strip gcc "" no-cfi-stripped

ls -la binaries/
rm -f /tmp/cfi.c
