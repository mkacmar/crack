#!/bin/sh
set -ex

ARCH=$1
SRC=test/e2e/testdata/main.c
mkdir -p binaries

. test/e2e/testdata/log-env.sh

RUNPATH_FLAGS="-Wl,--enable-new-dtags"

build() { $1 $RUNPATH_FLAGS -Wl,-rpath,$2 -o binaries/${ARCH}-$1-runpath-$3 $SRC; }

gcc -o binaries/${ARCH}-gcc-no-runpath $SRC
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
gcc $RUNPATH_FLAGS '-Wl,-rpath,$ORIGIN/../lib' -o binaries/${ARCH}-gcc-runpath-origin-relative $SRC

clang -o binaries/${ARCH}-clang-no-runpath $SRC
build clang /usr/lib absolute
build clang . dot
build clang /tmp tmp
clang -c -o binaries/${ARCH}-clang-relocatable.o $SRC

ls -la binaries/
