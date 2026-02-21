# Programmatic Usage

The following public packages can be used as a library:

- [`binary`](https://pkg.go.dev/go.kacmar.sk/crack/binary)
- [`rule`](https://pkg.go.dev/go.kacmar.sk/crack/rule)
- [`rule/elf`](https://pkg.go.dev/go.kacmar.sk/crack/rule/elf)
- [`rule/registry`](https://pkg.go.dev/go.kacmar.sk/crack/rule/registry)
- [`toolchain`](https://pkg.go.dev/go.kacmar.sk/crack/toolchain)

The standard workflow is to parse a binary, run rules, and inspect the findings:

```go
bin, _ := binary.ParseELF(f)

rules := registry.Where[rule.ELFRule](nil)

findings := rule.Check(rules, bin.Info, func(r rule.ELFRule) rule.Result {
    return r.Execute(bin)
})
```

[`Check`](https://pkg.go.dev/go.kacmar.sk/crack/rule#Check) handles applicability automatically — rules that don't match the binary's platform or compiler are skipped. See the [rules reference](rules.md) for available built-in rules.

For filtering, applicability checks, and other utilities, see the [package documentation](https://pkg.go.dev/go.kacmar.sk/crack).

## Custom Rules

To create a custom rule, implement the [`ELFRule`](https://pkg.go.dev/go.kacmar.sk/crack/rule#ELFRule) interface. For example, a rule that checks for a minimum stack size:

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

To detect custom compilers, implement [`ELFDetector`](https://pkg.go.dev/go.kacmar.sk/crack/toolchain#ELFDetector) and pass it to [`ParseELFWithDetector()`](https://pkg.go.dev/go.kacmar.sk/crack/binary#ParseELFWithDetector). This enables applicability checks for binaries built with internal or proprietary compilers:

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

For complete API documentation, see [pkg.go.dev](https://pkg.go.dev/go.kacmar.sk/crack).
