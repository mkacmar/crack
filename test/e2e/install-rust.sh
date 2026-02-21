#!/bin/sh
set -e

if ! command -v rustup >/dev/null 2>&1; then
    curl --proto '=https' --tlsv1.2 -sSf https://sh.rustup.rs | sh -s -- -y
    . "$HOME/.cargo/env"
fi

rustup update stable
rustup update nightly
