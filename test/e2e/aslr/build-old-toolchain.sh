#!/bin/sh
set -ex

SRC=test/e2e/testdata/main.c
mkdir -p binaries

gcc -fPIE -pie -Wl,-z,noexecstack -o binaries/x86_64-gcc-old-pie $SRC

ls -la binaries/
