#!/bin/sh
set -ex

C_SRC=test/e2e/testdata/main.c
RUST_SRC=test/e2e/testdata/main.rs
mkdir -p binaries

. test/e2e/testdata/log-env.sh

ARCH=$(uname -m)
if [ "$ARCH" != "x86_64" ]; then
    echo "Error: Retpoline is only supported on x86_64, detected $ARCH"
    exit 1
fi

gcc -mindirect-branch=thunk -mfunction-return=thunk -fcf-protection=none -o binaries/gcc-retpoline $C_SRC
gcc -fcf-protection=none -o binaries/gcc-no-retpoline $C_SRC

gcc -mindirect-branch=thunk -mfunction-return=thunk -fcf-protection=none -o binaries/gcc-retpoline-stripped $C_SRC
strip binaries/gcc-retpoline-stripped

clang -mretpoline -fcf-protection=none -o binaries/clang-retpoline $C_SRC
clang -fcf-protection=none -o binaries/clang-no-retpoline $C_SRC

gcc -fcf-protection=branch -o binaries/gcc-cet-ibt $C_SRC

rustc -o binaries/rustc-no-retpoline $RUST_SRC

ls -la binaries/
