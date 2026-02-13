# CRACK - Compiler Hardening Checker

> **Note**: This is a v0 release, API may change.

A tool to analyze ELF binaries for security hardening features.
Supports binaries compiled with `gcc`, `clang`, and `rustc` (stable).

Based on recommendations from:
- [OpenSSF Compiler Options Hardening Guide](https://best.openssf.org/Compiler-Hardening-Guides/Compiler-Options-Hardening-Guide-for-C-and-C++.html)
- [Gentoo Hardened Toolchain](https://wiki.gentoo.org/wiki/Hardened/Toolchain)
- [Debian Hardening](https://wiki.debian.org/Hardening)


## Installation

```sh
go install github.com/mkacmar/crack/cmd/crack@latest
```

Or download pre-built binaries from [releases](https://github.com/mkacmar/crack/releases).

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

See [rules reference](docs/rules.md) for available rules.

- `--rules <ids>` - Comma-separated list of rule IDs to run
- `--target-compiler <spec>` - Only run rules available for these compilers (e.g., `gcc`, `clang:15`)
- `--target-platform <spec>` - Only run rules available for these platforms (e.g., `arm64`, `amd64`)

The `--target-compiler` and `--target-platform` flags filter which rules are loaded based on their applicability.
At runtime, the tool also detects the actual compiler from binary metadata and skips rules that don't apply to the detected compiler.
For stripped binaries where detection fails, all loaded rules run.

### Output Options

- `--include-passed` - Include passing checks in output
- `--include-skipped` - Include skipped checks in output
- `--sarif <file>` - Save detailed SARIF report to file
- `--aggregate` - Aggregate findings into actionable recommendations
- `--exit-zero` - Exit with 0 even when findings are detected

The `--include-passed` and `--include-skipped` flags affect both text and SARIF output.

For programmatic access to results, use SARIF output (`--sarif`). [SARIF](https://sarifweb.azurewebsites.net/) (Static Analysis Results Interchange Format) is a standardized JSON format. We support SARIF version 2.1.0.

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

### Exit Codes

- `0` - Success (no findings, or `--exit-zero` specified)
- `1` - Error (invalid arguments, file errors, etc.)
- `2` - Findings detected


## Programmatic Usage

The public packages can be used as a library. See [API documentation](docs/api.md) for details on parsing binaries, running rules, writing custom rules, and custom compiler detection.


## License

MIT License - see [LICENSE](LICENSE) for details.
