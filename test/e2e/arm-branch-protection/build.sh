#!/bin/sh
set -ex

SRC=test/e2e/testdata/main.c
mkdir -p binaries

. test/e2e/testdata/log-env.sh

ARCH=$(uname -m)
if [ "$ARCH" != "aarch64" ]; then
    echo "Error: ARM branch protection is only supported on aarch64, detected $ARCH"
    exit 1
fi

gcc -mbranch-protection=standard -o binaries/gcc-branch-protection-standard $SRC
gcc -mbranch-protection=pac-ret -o binaries/gcc-branch-protection-pac-ret $SRC
gcc -mbranch-protection=bti -o binaries/gcc-branch-protection-bti $SRC
gcc -mbranch-protection=none -o binaries/gcc-no-branch-protection $SRC
gcc -mbranch-protection=standard -o binaries/gcc-branch-protection-stripped $SRC
strip binaries/gcc-branch-protection-stripped

clang -mbranch-protection=standard -o binaries/clang-branch-protection-standard $SRC
clang -mbranch-protection=pac-ret -o binaries/clang-branch-protection-pac-ret $SRC
clang -mbranch-protection=bti -o binaries/clang-branch-protection-bti $SRC
clang -mbranch-protection=none -o binaries/clang-no-branch-protection $SRC
clang -mbranch-protection=standard -o binaries/clang-branch-protection-stripped $SRC
strip binaries/clang-branch-protection-stripped

ls -la binaries/
