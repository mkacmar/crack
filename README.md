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
crack analyze [options] [<path>...]
```

### Input Options

- `<path>...` - Files or directories to analyze (supports glob patterns)
- `--recursive` - Recursively scan directories
- `--input <file>` - Read paths from file, one per line (use `-` for stdin)
- `--parallel <n>` - Number of files to analyze in parallel

### Rule Selection

See [Rules Reference](https://github.com/mkacmar/crack/wiki/Rules) for available rules.

- `--rules <ids>` - Comma-separated list of rule IDs to run
- `--target-compiler <spec>` - Only run rules available for these compilers (e.g., `gcc`, `clang:15`)
- `--target-platform <spec>` - Only run rules available for these platforms (e.g., `arm64`, `amd64`)

### Output Options

- `--include-passed` - Include passing checks in output
- `--include-skipped` - Include skipped checks in output
- `--sarif <file>` - Save detailed SARIF report to file
- `--aggregate` - Aggregate findings into actionable recommendations
- `--exit-zero` - Exit with 0 even when findings are detected

### Logging Options

- `--log <file>` - Write logs to file
- `--log-level <level>` - Log level: `none`, `debug`, `info`, `warn`, `error`

### Debuginfod Options

Fetch debug symbols from [debuginfod](https://sourceware.org/elfutils/Debuginfod.html) servers.

- `--debuginfod` - Enable debuginfod integration
- `--debuginfod-servers <urls>` - Comma-separated server URLs
- `--debuginfod-cache <dir>` - Cache directory for downloaded symbols
- `--debuginfod-timeout <duration>` - HTTP timeout
- `--debuginfod-retries <n>` - Max retries per server

## Documentation

- [Rules Reference](https://github.com/mkacmar/crack/wiki/Rules)

## License

MIT License - see [LICENSE](LICENSE) for details.
