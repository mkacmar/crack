#!/bin/sh
set -ex

ARCH=$1
SRC=test/e2e/testdata/main.c
mkdir -p binaries

. test/e2e/testdata/log-env.sh

# use --disable-new-dtags to force RPATH instead of RUNPATH
RPATH_FLAGS="-Wl,--disable-new-dtags"

gcc -o binaries/${ARCH}-gcc-no-rpath $SRC
gcc $RPATH_FLAGS -Wl,-rpath,/usr/lib -o binaries/${ARCH}-gcc-rpath-absolute $SRC
gcc $RPATH_FLAGS -Wl,-rpath,/usr/lib:/usr/local/lib -o binaries/${ARCH}-gcc-rpath-multiple-absolute $SRC
gcc $RPATH_FLAGS -Wl,-rpath,. -o binaries/${ARCH}-gcc-rpath-dot $SRC
gcc $RPATH_FLAGS -Wl,-rpath,.. -o binaries/${ARCH}-gcc-rpath-dotdot $SRC
gcc $RPATH_FLAGS -Wl,-rpath,./lib -o binaries/${ARCH}-gcc-rpath-relative $SRC
gcc $RPATH_FLAGS -Wl,-rpath,../lib -o binaries/${ARCH}-gcc-rpath-parent-relative $SRC
gcc $RPATH_FLAGS -Wl,-rpath,/tmp -o binaries/${ARCH}-gcc-rpath-tmp $SRC
gcc $RPATH_FLAGS -Wl,-rpath,/var/tmp -o binaries/${ARCH}-gcc-rpath-var-tmp $SRC
gcc $RPATH_FLAGS -Wl,-rpath,/tmp/mylibs -o binaries/${ARCH}-gcc-rpath-tmp-subdir $SRC
gcc $RPATH_FLAGS -Wl,-rpath,/usr/lib::/usr/local/lib -o binaries/${ARCH}-gcc-rpath-empty-component $SRC
gcc $RPATH_FLAGS -Wl,-rpath,/usr/lib:. -o binaries/${ARCH}-gcc-rpath-mixed $SRC

clang -o binaries/${ARCH}-clang-no-rpath $SRC
clang $RPATH_FLAGS -Wl,-rpath,/usr/lib -o binaries/${ARCH}-clang-rpath-absolute $SRC
clang $RPATH_FLAGS -Wl,-rpath,. -o binaries/${ARCH}-clang-rpath-dot $SRC
clang $RPATH_FLAGS -Wl,-rpath,/tmp -o binaries/${ARCH}-clang-rpath-tmp $SRC
clang -c -o binaries/${ARCH}-clang-relocatable.o $SRC

ls -la binaries/
