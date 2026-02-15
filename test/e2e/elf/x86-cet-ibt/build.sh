#!/bin/sh
set -ex

C_SRC=test/e2e/elf/testdata/main.c
RUST_SRC=test/e2e/elf/testdata/main.rs
mkdir -p binaries

. test/e2e/elf/testdata/log-env.sh

ARCH=$(uname -m)
if [ "$ARCH" != "amd64" ]; then
    echo "Error: Intel CET is only supported on amd64, detected $ARCH"
    exit 1
fi

build_c() { $1 -fcf-protection=$2 -o binaries/$1-cet-$2 $C_SRC; }
build_c_strip() { $1 -fcf-protection=$2 -o binaries/$1-cet-$3 $C_SRC && strip binaries/$1-cet-$3; }

build_c gcc full
build_c gcc branch
build_c gcc none
build_c_strip gcc full full-stripped

build_c clang full
build_c clang branch
build_c clang none

rustc -o binaries/rustc-default $RUST_SRC

ls -la binaries/
