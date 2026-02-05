#!/bin/sh
set -ex

SRC=test/e2e/testdata/main.c
mkdir -p binaries

. test/e2e/testdata/log-env.sh

ARCH=$(uname -m)
if [ "$ARCH" != "x86_64" ]; then
    echo "Error: Intel CET is only supported on x86_64, detected $ARCH"
    exit 1
fi

build() { $1 -fcf-protection=$2 -o binaries/$1-cet-$2 $SRC; }
build_strip() { $1 -fcf-protection=$2 -o binaries/$1-cet-$3 $SRC && strip binaries/$1-cet-$3; }

build gcc full
build gcc branch
build gcc none
build_strip gcc full full-stripped

build clang full
build clang branch
build clang none

ls -la binaries/
