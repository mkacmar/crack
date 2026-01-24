#!/bin/sh
set -ex

SRC=test/e2e/testdata/main.c
mkdir -p binaries

. test/e2e/testdata/log-env.sh

ARCH=$(uname -m)
if [ "$ARCH" != "x86_64" ]; then
    echo "Error: Retpoline is only supported on x86_64, detected $ARCH"
    exit 1
fi

gcc -mindirect-branch=thunk -mfunction-return=thunk -fcf-protection=none -o binaries/gcc-retpoline $SRC
gcc -fcf-protection=none -o binaries/gcc-no-retpoline $SRC

gcc -mindirect-branch=thunk -mfunction-return=thunk -fcf-protection=none -o binaries/gcc-retpoline-stripped $SRC
strip binaries/gcc-retpoline-stripped

clang -mretpoline -fcf-protection=none -o binaries/clang-retpoline $SRC
clang -fcf-protection=none -o binaries/clang-no-retpoline $SRC

gcc -fcf-protection=branch -o binaries/gcc-cet-ibt $SRC

ls -la binaries/
