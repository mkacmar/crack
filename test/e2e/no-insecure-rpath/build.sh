#!/bin/sh
set -ex

ARCH=$1
SRC=test/e2e/testdata/main.c
mkdir -p binaries

. test/e2e/testdata/log-env.sh

RPATH_FLAGS="-Wl,--disable-new-dtags"

build() { $1 $RPATH_FLAGS -Wl,-rpath,$2 -o binaries/${ARCH}-$1-rpath-$3 $SRC; }

gcc -o binaries/${ARCH}-gcc-no-rpath $SRC
build gcc /usr/lib absolute
build gcc /usr/lib:/usr/local/lib multiple-absolute
build gcc . dot
build gcc .. dotdot
build gcc ./lib relative
build gcc ../lib parent-relative
build gcc /tmp tmp
build gcc /var/tmp var-tmp
build gcc /tmp/mylibs tmp-subdir
build gcc /usr/lib::/usr/local/lib empty-component
build gcc /usr/lib:. mixed
build gcc lib bare-relative
build gcc subdir/lib subdir-relative
build gcc /dev/shm dev-shm
gcc $RPATH_FLAGS '-Wl,-rpath,$ORIGIN/../lib' -o binaries/${ARCH}-gcc-rpath-origin-relative $SRC

clang -o binaries/${ARCH}-clang-no-rpath $SRC
build clang /usr/lib absolute
build clang . dot
build clang /tmp tmp
clang -c -o binaries/${ARCH}-clang-relocatable.o $SRC

ls -la binaries/
