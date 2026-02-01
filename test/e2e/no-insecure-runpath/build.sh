#!/bin/sh
set -ex

ARCH=$1
SRC=test/e2e/testdata/main.c
mkdir -p binaries

. test/e2e/testdata/log-env.sh

# use --enable-new-dtags to force RUNPATH (default on modern linkers, but explicit for clarity)
RUNPATH_FLAGS="-Wl,--enable-new-dtags"

gcc -o binaries/${ARCH}-gcc-no-runpath $SRC
gcc $RUNPATH_FLAGS -Wl,-rpath,/usr/lib -o binaries/${ARCH}-gcc-runpath-absolute $SRC
gcc $RUNPATH_FLAGS -Wl,-rpath,/usr/lib:/usr/local/lib -o binaries/${ARCH}-gcc-runpath-multiple-absolute $SRC
gcc $RUNPATH_FLAGS -Wl,-rpath,. -o binaries/${ARCH}-gcc-runpath-dot $SRC
gcc $RUNPATH_FLAGS -Wl,-rpath,.. -o binaries/${ARCH}-gcc-runpath-dotdot $SRC
gcc $RUNPATH_FLAGS -Wl,-rpath,./lib -o binaries/${ARCH}-gcc-runpath-relative $SRC
gcc $RUNPATH_FLAGS -Wl,-rpath,../lib -o binaries/${ARCH}-gcc-runpath-parent-relative $SRC
gcc $RUNPATH_FLAGS -Wl,-rpath,/tmp -o binaries/${ARCH}-gcc-runpath-tmp $SRC
gcc $RUNPATH_FLAGS -Wl,-rpath,/var/tmp -o binaries/${ARCH}-gcc-runpath-var-tmp $SRC
gcc $RUNPATH_FLAGS -Wl,-rpath,/tmp/mylibs -o binaries/${ARCH}-gcc-runpath-tmp-subdir $SRC
gcc $RUNPATH_FLAGS -Wl,-rpath,/usr/lib::/usr/local/lib -o binaries/${ARCH}-gcc-runpath-empty-component $SRC
gcc $RUNPATH_FLAGS -Wl,-rpath,/usr/lib:. -o binaries/${ARCH}-gcc-runpath-mixed $SRC

clang -o binaries/${ARCH}-clang-no-runpath $SRC
clang $RUNPATH_FLAGS -Wl,-rpath,/usr/lib -o binaries/${ARCH}-clang-runpath-absolute $SRC
clang $RUNPATH_FLAGS -Wl,-rpath,. -o binaries/${ARCH}-clang-runpath-dot $SRC
clang $RUNPATH_FLAGS -Wl,-rpath,/tmp -o binaries/${ARCH}-clang-runpath-tmp $SRC
clang -c -o binaries/${ARCH}-clang-relocatable.o $SRC

ls -la binaries/
