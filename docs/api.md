# Programmatic Usage

The following public packages can be used as a library:

- [`binary`](https://pkg.go.dev/go.kacmar.sk/crack/binary)
- [`binary/elf`](https://pkg.go.dev/go.kacmar.sk/crack/binary/elf)
- [`rule`](https://pkg.go.dev/go.kacmar.sk/crack/rule)
- [`rule/elf`](https://pkg.go.dev/go.kacmar.sk/crack/rule/elf)
- [`rule/registry`](https://pkg.go.dev/go.kacmar.sk/crack/rule/registry)
- [`toolchain`](https://pkg.go.dev/go.kacmar.sk/crack/toolchain)

The standard workflow is to parse a binary, classify it, run rules, and inspect the findings:

```go
bin, _ := elf.Open(f)

profile := binary.Profile{
    Architecture: elf.DetectArchitecture(bin),
    LibC:         elf.DetectLibC(bin),
    Toolchain:    elf.DefaultToolchainDetector{}.Detect(bin),
}

rules := registry.Where[rule.ELFRule](nil)

findings := rule.Check(rules, profile, func(r rule.ELFRule) rule.Result {
    return r.Execute(bin)
})
```

[`Check`](https://pkg.go.dev/go.kacmar.sk/crack/rule#Check) handles applicability automatically - rules that don't match the binary's platform or compiler are skipped. See the [rules reference](rules.md) for available built-in rules.

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

func (r MinStackSizeRule) Execute(bin elf.Binary) rule.Result {
    for _, prog := range bin.Progs() {
        if prog.Type == stdelf.PT_GNU_STACK && prog.Memsz >= r.MinBytes {
            return rule.Result{Status: rule.StatusPassed, Message: fmt.Sprintf("Stack size %d bytes", prog.Memsz)}
        }
    }
    return rule.Result{Status: rule.StatusFailed, Message: "Stack size below minimum or not set"}
}
```

## Custom Compiler Detection

Toolchain detection has two extension points. To recognize a private compiler that signs `.comment` or DWARF `DW_AT_producer`, implement [`StringDetector`](https://pkg.go.dev/go.kacmar.sk/crack/toolchain#StringDetector) and set it on [`DefaultToolchainDetector`](https://pkg.go.dev/go.kacmar.sk/crack/binary/elf#DefaultToolchainDetector):

```go
// AcmeStringDetector recognizes Acme Corp's internal compiler, falling back to standard detection.
type AcmeStringDetector struct {
    fallback toolchain.DefaultStringDetector
}

func (d AcmeStringDetector) Detect(s string) (toolchain.Compiler, toolchain.Version) {
    if strings.Contains(s, "ACME C Compiler") {
        parts := strings.Fields(s)
        if len(parts) >= 4 {
            if v, err := toolchain.ParseVersion(parts[3]); err == nil {
                return toolchain.Compiler("acme-cc"), v
            }
        }
        return toolchain.Compiler("acme-cc"), toolchain.Version{}
    }
    return d.fallback.Detect(s)
}

detector := elf.DefaultToolchainDetector{StringDetector: AcmeStringDetector{}}
tc := detector.Detect(bin)
```

To use a different evidence source entirely (e.g. recognize a compiler by a dedicated section, similar to how Go is recognized via `.go.buildinfo`), implement [`ToolchainDetector`](https://pkg.go.dev/go.kacmar.sk/crack/binary/elf#ToolchainDetector) directly:

```go
type AcmeToolchainDetector struct {
    fallback elf.ToolchainDetector
}

func (d AcmeToolchainDetector) Detect(b elf.Binary) toolchain.Toolchain {
    if _, err := elf.FindSection(b, ".acme.buildinfo"); err == nil {
        return toolchain.Toolchain{Compiler: toolchain.Compiler("acme-cc")}
    }
    return d.fallback.Detect(b)
}

detector := AcmeToolchainDetector{fallback: elf.DefaultToolchainDetector{}}
tc := detector.Detect(bin)
```

## Custom Section Resolver

For stripped distro binaries the typical path is the built-in debuginfod integration via the CLI's `--debuginfod` flag, which transparently fetches missing sections by build ID from public servers like `debuginfod.fedoraproject.org`. Implement [`Resolver`](https://pkg.go.dev/go.kacmar.sk/crack/binary/elf#Resolver) only when sections need to come from a non-standard source - a corporate symbol store, an artifact registry, a CI cache - and wire it via [`WithResolverFactory`](https://pkg.go.dev/go.kacmar.sk/crack/binary/elf#WithResolverFactory):

```go
// ArtifactStoreResolver fetches stripped sections from an internal HTTP debug archive.
type ArtifactStoreResolver struct {
    baseURL string
    buildID string
    client  *http.Client
}

func (r ArtifactStoreResolver) FetchSection(name string) ([]byte, error) {
    url := fmt.Sprintf("%s/debug/%s/section/%s", r.baseURL, r.buildID, name)
    resp, err := r.client.Get(url)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()
    if resp.StatusCode != http.StatusOK {
        return nil, fmt.Errorf("section %s: %s", name, resp.Status)
    }
    return io.ReadAll(resp.Body)
}

factory := func(buildID string) elf.Resolver {
    return ArtifactStoreResolver{baseURL: "https://debug.internal", buildID: buildID, client: http.DefaultClient}
}
bin, _ := elf.Open(f, elf.WithResolverFactory(factory))
```

The factory receives the binary's build ID so the resolver can target the matching debug artifact. Sections fetched by the resolver are merged transparently into the [`Binary`](https://pkg.go.dev/go.kacmar.sk/crack/binary/elf#Binary) view, so rules see them as if they were on disk.

## Custom Binary Source

The [`Binary`](https://pkg.go.dev/go.kacmar.sk/crack/binary/elf#Binary) interface is the analysis input contract. Most extensions don't need a full implementation - embedding and overriding a single method is usually enough. For example, injecting a build ID from an external manifest for binaries that lack `NT_GNU_BUILD_ID`:

```go
type ManifestBuildID struct {
    elf.Binary
    id string
}

func (m ManifestBuildID) BuildID() string { return m.id }

bin, _ := elf.Open(f)
wrapped := ManifestBuildID{Binary: bin, id: manifest.BuildHash}
```

The wrapped binary plugs into the rest of the API unchanged. See the godoc for the full interface surface when a more comprehensive override is needed.

For complete API documentation, see [pkg.go.dev](https://pkg.go.dev/go.kacmar.sk/crack).
