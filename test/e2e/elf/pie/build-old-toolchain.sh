#!/bin/sh
set -ex

SRC=test/e2e/elf/testdata/main.c
mkdir -p binaries

gcc -fPIE -pie -o binaries/x86_64-gcc-old-pie $SRC

ls -la binaries/

