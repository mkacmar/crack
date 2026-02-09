#!/bin/sh
set -ex

ARCH=$1
C_SRC=test/e2e/testdata/main.c
RUST_SRC=test/e2e/testdata/main.rs
mkdir -p binaries

. test/e2e/testdata/log-env.sh

RUNPATH_FLAGS="-Wl,--enable-new-dtags"

build_c() { $1 $RUNPATH_FLAGS -Wl,-rpath,$2 -o binaries/${ARCH}-$1-runpath-$3 $C_SRC; }

gcc -o binaries/${ARCH}-gcc-no-runpath $C_SRC
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
gcc $RUNPATH_FLAGS '-Wl,-rpath,$ORIGIN/../lib' -o binaries/${ARCH}-gcc-runpath-origin-relative $C_SRC

clang -o binaries/${ARCH}-clang-no-runpath $C_SRC
build_c clang /usr/lib absolute
build_c clang . dot
build_c clang /tmp tmp
clang -c -o binaries/${ARCH}-clang-relocatable.o $C_SRC

rustc -o binaries/${ARCH}-rustc-no-runpath $RUST_SRC
rustc -C link-arg=--enable-new-dtags -C link-arg=-rpath -C link-arg=/usr/lib -o binaries/${ARCH}-rustc-runpath-absolute $RUST_SRC
rustc -C link-arg=--enable-new-dtags -C link-arg=-rpath -C link-arg=. -o binaries/${ARCH}-rustc-runpath-dot $RUST_SRC
rustc -C link-arg=--enable-new-dtags -C link-arg=-rpath -C link-arg=/tmp -o binaries/${ARCH}-rustc-runpath-tmp $RUST_SRC

ls -la binaries/
