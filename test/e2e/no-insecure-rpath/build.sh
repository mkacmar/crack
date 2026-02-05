#!/bin/sh
set -ex

ARCH=$1
SRC=test/e2e/testdata/main.c
mkdir -p binaries

. test/e2e/testdata/log-env.sh

RPATH_FLAGS="-Wl,--disable-new-dtags"

build_c() { $1 $RPATH_FLAGS -Wl,-rpath,$2 -o binaries/${ARCH}-$1-rpath-$3 $SRC; }

gcc -o binaries/${ARCH}-gcc-no-rpath $SRC
build_c gcc /usr/lib absolute
build_c gcc /usr/lib:/usr/local/lib multiple-absolute
build_c gcc . dot
build_c gcc .. dotdot
build_c gcc ./lib relative
build_c gcc ../lib parent-relative
build_c gcc /tmp tmp
build_c gcc /var/tmp var-tmp
build_c gcc /tmp/mylibs tmp-subdir
build_c gcc /usr/lib::/usr/local/lib empty-component
build_c gcc /usr/lib:. mixed
build_c gcc lib bare-relative
build_c gcc subdir/lib subdir-relative
build_c gcc /dev/shm dev-shm
gcc $RPATH_FLAGS '-Wl,-rpath,$ORIGIN/../lib' -o binaries/${ARCH}-gcc-rpath-origin-relative $SRC

clang -o binaries/${ARCH}-clang-no-rpath $SRC
build_c clang /usr/lib absolute
build_c clang . dot
build_c clang /tmp tmp
clang -c -o binaries/${ARCH}-clang-relocatable.o $SRC

ls -la binaries/
