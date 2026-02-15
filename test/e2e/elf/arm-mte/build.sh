#!/bin/sh
set -ex

NDK_DIR=$1
if [ -z "$NDK_DIR" ] || [ ! -d "$NDK_DIR" ]; then
    echo "Usage: $0 <ndk-dir>"
    exit 1
fi

C_SRC=test/e2e/elf/testdata/main.c
mkdir -p binaries

CLANG=${NDK_DIR}/toolchains/llvm/prebuilt/linux-amd64/bin/clang
STRIP=${NDK_DIR}/toolchains/llvm/prebuilt/linux-amd64/bin/llvm-strip
TARGET=arm64-linux-android35


$CLANG --version

$CLANG --target=$TARGET -march=armv8.5-a+memtag -fsanitize=memtag-stack,memtag-heap -o binaries/clang-mte $C_SRC
$CLANG --target=$TARGET -march=armv8.5-a+memtag -fsanitize=memtag-stack,memtag-heap -o binaries/clang-mte-stripped $C_SRC
$STRIP binaries/clang-mte-stripped

$CLANG --target=$TARGET -o binaries/clang-no-mte $C_SRC

ls -la binaries/
