#!/bin/sh
echo "=== Build environment ==="
uname -m
gcc --version | head -1
clang --version | head -1

