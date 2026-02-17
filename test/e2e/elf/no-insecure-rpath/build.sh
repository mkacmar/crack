#!/bin/sh
set -ex

ARCH=$1
C_SRC=test/e2e/elf/testdata/main.c
RUST_SRC=test/e2e/elf/testdata/main.rs
mkdir -p binaries

. test/e2e/elf/testdata/log-env.sh

RPATH_FLAGS="-Wl,--disable-new-dtags"

build_c() { $1 $RPATH_FLAGS -Wl,-rpath,$2 -o binaries/${ARCH}-$1-rpath-$3 $C_SRC; }

gcc -o binaries/${ARCH}-gcc-no-rpath $C_SRC
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
gcc $RPATH_FLAGS '-Wl,-rpath,$ORIGIN/../lib' -o binaries/${ARCH}-gcc-rpath-origin-relative $C_SRC
gcc $RPATH_FLAGS '-Wl,-rpath,$ORIGIN/..' -o binaries/${ARCH}-gcc-rpath-origin-parent $C_SRC

clang -o binaries/${ARCH}-clang-no-rpath $C_SRC
build_c clang /usr/lib absolute
build_c clang . dot
build_c clang /tmp tmp
clang -c -o binaries/${ARCH}-clang-relocatable.o $C_SRC

rustc -o binaries/${ARCH}-rustc-no-rpath $RUST_SRC
rustc -C link-arg=-Wl,--disable-new-dtags -C link-arg=-Wl,-rpath,/usr/lib -o binaries/${ARCH}-rustc-rpath-absolute $RUST_SRC
rustc -C link-arg=-Wl,--disable-new-dtags -C link-arg=-Wl,-rpath,. -o binaries/${ARCH}-rustc-rpath-dot $RUST_SRC
rustc -C link-arg=-Wl,--disable-new-dtags -C link-arg=-Wl,-rpath,/tmp -o binaries/${ARCH}-rustc-rpath-tmp $RUST_SRC

ls -la binaries/
