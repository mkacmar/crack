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

build_c() { $1 -mbranch-protection=$2 -o binaries/$1-$3 $SRC; }
build_c_strip() { $1 -mbranch-protection=$2 -o binaries/$1-$3 $SRC && strip binaries/$1-$3; }

build_c gcc standard branch-protection-standard
build_c gcc pac-ret branch-protection-pac-ret
build_c gcc bti branch-protection-bti
build_c gcc none no-branch-protection
build_c_strip gcc standard branch-protection-stripped

build_c clang standard branch-protection-standard
build_c clang pac-ret branch-protection-pac-ret
build_c clang bti branch-protection-bti
build_c clang none no-branch-protection
build_c_strip clang standard branch-protection-stripped

ls -la binaries/
