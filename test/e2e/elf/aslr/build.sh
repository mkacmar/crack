#!/bin/sh
set -ex

ARCH=$1
C_SRC=test/e2e/elf/testdata/main.c
mkdir -p binaries

. test/e2e/elf/testdata/log-env.sh

build_c() { $1 $2 -o binaries/${ARCH}-$1-$3 $C_SRC; }
build_c_strip() { $1 $2 -o binaries/${ARCH}-$1-$3 $C_SRC && strip binaries/${ARCH}-$1-$3; }

build_c gcc "-fPIE -pie -Wl,-z,noexecstack" aslr-full
build_c gcc "-fno-pie -no-pie" no-pie
build_c gcc "-fPIE -pie -Wl,-z,execstack" execstack
build_c gcc "-shared -fPIC" shared
build_c gcc "-static-pie -Wl,-z,noexecstack" static-pie
build_c_strip gcc "-fPIE -pie -Wl,-z,noexecstack" aslr-stripped
build_c gcc "-static -fno-pie -no-pie" static-no-pie

build_c clang "-fPIE -pie -Wl,-z,noexecstack" aslr-full
build_c clang "-fno-pie -no-pie" no-pie
build_c clang "-fPIE -pie -Wl,-z,execstack" execstack
build_c clang "-static -fno-pie -no-pie" static-no-pie

gcc -fPIE -pie -Wl,-z,noexecstack -o binaries/${ARCH}-gcc-textrel-patched $C_SRC
# Patch DT_DEBUG -> DT_TEXTREL in the dynamic section.
gcc -xc -o /tmp/patch-textrel - <<'EOF'
#define _GNU_SOURCE
#include <string.h>
#include <stdio.h>
#include <stdlib.h>
int main(int argc, char **argv) {
    FILE *f = fopen(argv[1], "r+b");
    fseek(f, 0, SEEK_END);
    long sz = ftell(f);
    rewind(f);
    char *buf = malloc(sz);
    fread(buf, 1, sz, f);
    unsigned long entry[2] = {21, 0};
    char *p = memmem(buf, sz, entry, sizeof(entry));
    if (!p) { fputs("DT_DEBUG not found\n", stderr); return 1; }
    unsigned long tag = 22;
    memcpy(p, &tag, sizeof(tag));
    rewind(f);
    fwrite(buf, 1, sz, f);
    fclose(f);
    return 0;
}
EOF
/tmp/patch-textrel binaries/${ARCH}-gcc-textrel-patched

gcc -c -o binaries/${ARCH}-gcc-relocatable.o $C_SRC

ls -la binaries/
