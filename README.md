# CRACK - Compiler Hardening Checker

> **Work in Progress**: This project is under active development. Functionality may change without notice.

A tool to analyze ELF binaries for security hardening features.

Based on recommendations from:
- [OpenSSF Compiler Options Hardening Guide for C and C++](https://best.openssf.org/Compiler-Hardening-Guides/Compiler-Options-Hardening-Guide-for-C-and-C++.html)
- [Gentoo Hardened Toolchain](https://wiki.gentoo.org/wiki/Hardened/Toolchain)
- [Debian Hardening](https://wiki.debian.org/Hardening)
- [Ubuntu Toolchain Compiler Flags](https://wiki.ubuntu.com/ToolChain/CompilerFlags)

## Installation

```sh
go install github.com/mkacmar/crack/cmd/crack@latest
```

Or download pre-built binaries from [Releases](https://github.com/mkacmar/crack/releases).

## Usage

```sh
# Analyze a binary
crack analyze /usr/bin/ls
crack analyze /usr/bin/*

# Analyze a directory recursively
crack analyze --recursive /usr/bin/

# Analyze paths from file, one path per line (or stdin with "-")
crack analyze --input files.txt

# Analyze with specific rules, comma-separated (see wiki for rule IDs)
crack analyze --rules pie,relro,stack-canary /usr/bin/ls

# Analyze and generate SARIF output
crack analyze --sarif results.sarif /usr/bin/ls

# Analyze with debug symbols from debuginfod (comma-separated URLs)
crack analyze --debuginfod /usr/bin/ls
crack analyze --debuginfod --debuginfod-urls https://debuginfod.archlinux.org /usr/bin/ls
```

Use `crack analyze --help` for all options.

## Documentation

- [Rules Reference](https://github.com/mkacmar/crack/wiki/Rules)

## License

MIT License - see [LICENSE](LICENSE) for details.
