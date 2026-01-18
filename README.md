# CRACK - Compiler Hardening Checker

A tool to analyze ELF binaries for security hardening features.

Based on recommendations from:
- [OpenSSF Compiler Options Hardening Guide for C and C++](https://best.openssf.org/Compiler-Hardening-Guides/Compiler-Options-Hardening-Guide-for-C-and-C++.html)
- [Gentoo Hardened Toolchain](https://wiki.gentoo.org/wiki/Hardened/Toolchain)
- [Debian Hardening](https://wiki.debian.org/Hardening)
- [Ubuntu Toolchain Compiler Flags](https://wiki.ubuntu.com/ToolChain/CompilerFlags)

## Usage

```bash
# Show help
crack analyze --help

# Analyze a binary with the default (recommended) profile
crack analyze /usr/bin/ls

# List rules in a specific profile
crack analyze --profile=hardened --list-rules

# Analyze with debuginfod to fetch debug symbols for stripped binaries
crack analyze --profile=hardened --debuginfod --debuginfod-urls=https://debuginfod.elfutils.org /usr/bin/ls
```

## Available Rules

### Universal Rules

| Rule ID | Description | GCC | Clang (LLVM) | Compile time | Link time | Performance Impact |
|---------|-------------|-----|--------------|--------------|-----------|--------------------|
| `asan` | [AddressSanitizer](https://clang.llvm.org/docs/AddressSanitizer.html) | ✓ | ✓ | ✓ | | ✓ |
| `aslr` | [Address Space Layout Randomization compatibility](https://en.wikipedia.org/wiki/Address_space_layout_randomization) | ✓ | ✓ | ✓ | ✓ | |
| `cfi` | [Control Flow Integrity](https://clang.llvm.org/docs/ControlFlowIntegrity.html) | | ✓ | ✓ | | ✓ |
| `fortify-source` | [FORTIFY_SOURCE buffer overflow protection](https://gcc.gnu.org/onlinedocs/gcc/Instrumentation-Options.html#index-D_FORTIFY_SOURCE) | ✓ | ✓ | ✓ | | |
| `full-relro` | [Full RELRO (immediate binding)](https://www.redhat.com/en/blog/hardening-elf-binaries-using-relocation-read-only-relro) | ✓ | ✓ | | ✓ | |
| `hidden-symbols` | [Hidden symbol visibility](https://gcc.gnu.org/onlinedocs/gcc/Code-Gen-Options.html#index-fvisibility) | ✓ | ✓ | ✓ | | |
| `kernel-cfi` | [Kernel CFI](https://clang.llvm.org/docs/ControlFlowIntegrity.html#fsanitize-kcfi) | | ✓ | ✓ | | ✓ |
| `no-dlopen` | [Dynamic library loading disabled](https://sourceware.org/binutils/docs/ld/Options.html) | ✓ | ✓ | | ✓ | |
| `no-dump` | [Core dump disabled](https://sourceware.org/binutils/docs/ld/Options.html) | ✓ | ✓ | | ✓ | |
| `no-insecure-rpath` | [No insecure RPATH entries](https://man7.org/linux/man-pages/man8/ld.so.8.html) | ✓ | ✓ | | ✓ | |
| `no-insecure-runpath` | [No insecure RUNPATH entries](https://man7.org/linux/man-pages/man8/ld.so.8.html) | ✓ | ✓ | | ✓ | |
| `no-plt` | [No PLT (direct GOT access)](https://gcc.gnu.org/onlinedocs/gcc/Code-Gen-Options.html#index-fno-plt) | ✓ | ✓ | ✓ | | |
| `nx-bit` | [Non-executable stack (NX bit)](https://www.gnu.org/software/binutils/manual/ld.html) | ✓ | ✓ | | ✓ | |
| `pie` | [Position Independent Executable](https://gcc.gnu.org/onlinedocs/gcc/Code-Gen-Options.html#index-fPIE) | ✓ | ✓ | ✓ | ✓ | |
| `relro` | [Read-only relocations (RELRO)](https://www.redhat.com/en/blog/hardening-elf-binaries-using-relocation-read-only-relro) | ✓ | ✓ | | ✓ | |
| `safe-stack` | [SafeStack](https://clang.llvm.org/docs/SafeStack.html) | | ✓ | ✓ | | ✓ |
| `separate-code` | [Separate code and data segments](https://sourceware.org/binutils/docs/ld/Options.html) | ✓ | ✓ | | ✓ | |
| `stack-canary` | [Stack canary protection](https://gcc.gnu.org/onlinedocs/gcc/Instrumentation-Options.html#index-fstack-protector) | ✓ | ✓ | ✓ | | |
| `stack-limit` | [Stack size limit](https://sourceware.org/binutils/docs/ld/Options.html) | ✓ | ✓ | | ✓ | |
| `stripped` | [Debug symbols stripped](https://sourceware.org/binutils/docs/binutils/strip.html) | ✓ | ✓ | | ✓ | |
| `ubsan` | [UndefinedBehaviorSanitizer](https://clang.llvm.org/docs/UndefinedBehaviorSanitizer.html) | ✓ | ✓ | ✓ | | ✓ |
| `wxorx` | [Write XOR Execute policy](https://en.wikipedia.org/wiki/W%5EX) | ✓ | ✓ | | ✓ | |

### x86/x86-64 Specific Rules

| Rule ID | Description | GCC | Clang (LLVM) | Compile time | Link time | Performance Impact |
|---------|-------------|-----|--------------|--------------|-----------|--------------------|
| `intel-cet-ibt` | [Intel CET Indirect Branch Tracking](https://gcc.gnu.org/onlinedocs/gcc/Instrumentation-Options.html#index-fcf-protection) | ✓ | ✓ | ✓ | | |
| `intel-cet-shstk` | [Intel CET Shadow Stack](https://gcc.gnu.org/onlinedocs/gcc/Instrumentation-Options.html#index-fcf-protection) | ✓ | ✓ | ✓ | | |
| `x86-retpoline` | [Retpoline Spectre v2 mitigation](https://gcc.gnu.org/onlinedocs/gcc/x86-Options.html#index-mindirect-branch) | ✓ | ✓ | ✓ | | ✓ |

### ARM64 Specific Rules

| Rule ID | Description | GCC | Clang (LLVM) | Compile time | Link time | Performance Impact |
|---------|-------------|-----|--------------|--------------|-----------|--------------------|
| `arm-branch-protection` | [ARM branch protection](https://gcc.gnu.org/onlinedocs/gcc/AArch64-Options.html#index-mbranch-protection) | ✓ | ✓ | ✓ | | |
| `arm-bti` | [ARM Branch Target Identification](https://gcc.gnu.org/onlinedocs/gcc/AArch64-Options.html#index-mbranch-protection) | ✓ | ✓ | ✓ | | |
| `arm-mte` | [ARM Memory Tagging Extension](https://clang.llvm.org/docs/MemtagSanitizer.html) | | ✓ | ✓ | | ✓ |
| `arm-pac` | [ARM Pointer Authentication](https://gcc.gnu.org/onlinedocs/gcc/AArch64-Options.html#index-mbranch-protection) | ✓ | ✓ | ✓ | | |

## License

MIT License - see [LICENSE](LICENSE) for details.
