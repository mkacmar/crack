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

build() { $1 -mbranch-protection=$2 -o binaries/$1-$3 $SRC; }
build_strip() { $1 -mbranch-protection=$2 -o binaries/$1-$3 $SRC && strip binaries/$1-$3; }

build gcc standard branch-protection-standard
build gcc pac-ret branch-protection-pac-ret
build gcc bti branch-protection-bti
build gcc none no-branch-protection
build_strip gcc standard branch-protection-stripped

build clang standard branch-protection-standard
build clang pac-ret branch-protection-pac-ret
build clang bti branch-protection-bti
build clang none no-branch-protection
build_strip clang standard branch-protection-stripped

ls -la binaries/
