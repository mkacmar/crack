# Programmatic Usage

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

Run rules against the parsed binary (see [rules reference](rules.md) for available rules):

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

## Applicability

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

## Custom Rules

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

## Custom Compiler Detection

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
