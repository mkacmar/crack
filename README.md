# CRACK - Compiler Hardening Checker

> **Work in Progress**: This project is currently unstable and under active development. APIs, rules, and output formats may change without notice. Do not use in production workloads.

A tool to analyze ELF binaries for security hardening features.

Based on recommendations from:
- [OpenSSF Compiler Options Hardening Guide for C and C++](https://best.openssf.org/Compiler-Hardening-Guides/Compiler-Options-Hardening-Guide-for-C-and-C++.html)
- [Gentoo Hardened Toolchain](https://wiki.gentoo.org/wiki/Hardened/Toolchain)
- [Debian Hardening](https://wiki.debian.org/Hardening)
- [Ubuntu Toolchain Compiler Flags](https://wiki.ubuntu.com/ToolChain/CompilerFlags)

## Usage

```bash
crack analyze /usr/bin/ls
crack analyze --rules=pie,relro,stack-canary /usr/bin/ls
crack analyze --debuginfod /usr/bin/ls
```


## License

MIT License - see [LICENSE](LICENSE) for details.
