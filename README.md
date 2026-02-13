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

See [rules reference](https://github.com/mkacmar/crack/wiki/Rules) for available rules.

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

The public packages can be used as a library to integrate binary analysis into your own tools, write custom rules, or build an alternative frontend.

Import the relevant packages:

```go
import (
    "github.com/mkacmar/crack/binary"
    "github.com/mkacmar/crack/rule"
    "github.com/mkacmar/crack/rule/elf"
    "github.com/mkacmar/crack/toolchain"
)
```

Parse a binary using `binary.ParseELF`:

```go
f, err := os.Open("/usr/bin/ls")
if err != nil {
    log.Fatal(err)
}
defer f.Close()

bin, err := binary.ParseELF(f)
if err != nil {
    log.Fatal(err)
}
```

The `ELFBinary` struct provides access to ELF metadata (see [`debug/elf`](https://pkg.go.dev/debug/elf) for types), symbol tables, and detected toolchain info. Helper methods simplify common checks:

- `HasDynTag()`, `HasDynFlag()`, `DynString()` - query [dynamic section tags](https://man7.org/linux/man-pages/man5/elf.5.html) (`DT_*`)
- `HasGNUProperty()` - check [GNU program properties](https://docs.kernel.org/userspace-api/ELF.html) for features like CET, BTI

Run rules against the parsed binary (see [rules reference](https://github.com/mkacmar/crack/wiki/Rules) for available rules):

```go
rules := []rule.ELFRule{
    elf.PIERule{},
    elf.StackCanaryRule{},
    elf.FullRELRORule{},
}

findings := rule.Check(rules, bin.Info, func(r rule.ELFRule) rule.Result {
    return r.Execute(bin)
})
```

All rules are checked by default - passed, failed, and skipped findings are all included in results.

### Applicability

Each rule declares its applicability - which platforms and compilers it supports. For example, a rule might require GCC 10+ or only apply to ARM64 architecture. When analyzing a binary, rules that don't apply are automatically skipped.

Rules specify a `MinVersion` - the minimum compiler version required for the feature. The `Platform` field specifies which architectures the rule applies to. For example, the ARM PAC rule:

```go
rule.Applicability{
    Platform: binary.PlatformARM64v83,
    Compilers: map[toolchain.Compiler]rule.CompilerRequirement{
        toolchain.GCC:   {MinVersion: toolchain.Version{Major: 10, Minor: 1}, Flag: "-mbranch-protection=pac-ret"},
        toolchain.Clang: {MinVersion: toolchain.Version{Major: 12}, Flag: "-mbranch-protection=pac-ret"},
    },
}
```

If your binaries are built with an internal compiler, register it via a custom detector so rules can determine whether they apply.

Use `rule.CheckApplicability()` to manually check if a rule applies:

```go
result := rule.CheckApplicability(myRule.Applicability(), bin.Info)
if result == rule.Applicable {
	// ...
}
```

Use `rule.FilterRules()` to pre-filter rules based on your target environment. The filter uses `MaxVersion` to exclude rules that require a newer compiler than you have:

```go
filter := &rule.TargetFilter{
    Compilers: []rule.CompilerTarget{
        {Compiler: toolchain.GCC, MaxVersion: &toolchain.Version{Major: 12}},
    },
}
filtered := rule.FilterRules(rules, filter)
```

### Custom Rules

To create a custom rule, implement the `rule.ELFRule` interface. For example, a rule that checks for a minimum stack size:

```go
type MinStackSizeRule struct {
    MinBytes uint64
}

func (r MinStackSizeRule) ID() string          { return "min-stack-size" }
func (r MinStackSizeRule) Name() string        { return "Minimum Stack Size" }
func (r MinStackSizeRule) Description() string { return "Ensures stack size meets minimum requirements" }

func (r MinStackSizeRule) Applicability() rule.Applicability {
    return rule.Applicability{
        Platform: binary.PlatformAll,
    }
}

func (r MinStackSizeRule) Execute(bin *binary.ELFBinary) rule.Result {
    for _, prog := range bin.Progs {
        if prog.Type == elf.PT_GNU_STACK && prog.Memsz >= r.MinBytes {
            return rule.Result{Status: rule.StatusPassed, Message: fmt.Sprintf("Stack size %d bytes", prog.Memsz)}
        }
    }
    return rule.Result{Status: rule.StatusFailed, Message: "Stack size below minimum or not set"}
}
```

### Custom Compiler Detection

To detect custom compilers, implement `toolchain.ELFDetector` and pass it to `binary.ParseELFWithDetector()`. This enables applicability checks for binaries built with internal or proprietary compilers:

```go
// AcmeDetector detects Acme Corp's internal compiler, falling back to standard detection.
type AcmeDetector struct {
    fallback toolchain.ELFCommentDetector
}

func (d AcmeDetector) Detect(comment string) (toolchain.Compiler, toolchain.Version) {
    // Acme compiler writes "ACME C Compiler 2.3.1" in .comment section
    if strings.Contains(comment, "ACME C Compiler") {
        parts := strings.Fields(comment)
        if len(parts) >= 4 {
            if v, err := toolchain.ParseVersion(parts[3]); err == nil {
                return toolchain.Compiler("acme-cc"), v
            }
        }
        return toolchain.Compiler("acme-cc"), toolchain.Version{}
    }
    return d.fallback.Detect(comment)
}

bin, err := binary.ParseELFWithDetector(f, AcmeDetector{})
```

For complete API documentation, see [pkg.go.dev](https://pkg.go.dev/github.com/mkacmar/crack).


## License

MIT License - see [LICENSE](LICENSE) for details.
