#!/bin/sh
set -ex

ARCH=$1
SRC=test/e2e/testdata/plt.c
mkdir -p binaries

. test/e2e/testdata/log-env.sh

case "$ARCH" in
    x86_64)
        gcc -fno-plt -fcf-protection=full -o binaries/gcc-no-plt-cet $SRC
        gcc -fno-plt -fcf-protection=none -o binaries/gcc-no-plt $SRC
        gcc -fno-plt -fcf-protection=none -o binaries/gcc-no-plt-stripped $SRC
        strip binaries/gcc-no-plt-stripped
        clang -fno-plt -o binaries/clang-no-plt $SRC
        clang -fno-plt -o binaries/clang-no-plt-stripped $SRC
        strip binaries/clang-no-plt-stripped
        gcc -fcf-protection=none -o binaries/gcc-plt $SRC
        clang -o binaries/clang-plt $SRC
        ;;
    i386)
        gcc -m32 -fno-plt -fcf-protection=full -o binaries/i386-gcc-no-plt-cet $SRC
        gcc -m32 -fno-plt -fcf-protection=none -o binaries/i386-gcc-no-plt $SRC
        gcc -m32 -fcf-protection=none -o binaries/i386-gcc-plt $SRC
        clang -m32 -fno-plt -o binaries/i386-clang-no-plt $SRC
        clang -m32 -o binaries/i386-clang-plt $SRC
        ;;
    *)
        echo "no-plt rule only supports x86_64 and i386"
        exit 1
        ;;
esac

ls -la binaries/
