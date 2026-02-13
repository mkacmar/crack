## ARM Branch Protection

- **Rule ID:** `arm-branch-protection`
- **Implementation:** `ARMBranchProtectionRule`

Checks for ARM branch protection (BTI + PAC combined). This enables both Branch Target Identification to validate indirect branch targets and Pointer Authentication to sign return addresses.

### Platform

arm64 (requires ISA v8.5+)

### Toolchain

| Compiler | Minimal Version | Default Version | Flag |
|:---------|:----------------|:----------------|:-----|
| clang | 12.0 | - | `-mbranch-protection=standard` |
| gcc | 10.1 | - | `-mbranch-protection=standard` |


---

## ARM Branch Target Identification

- **Rule ID:** `arm-bti`
- **Implementation:** `ARMBTIRule`

Checks for ARM Branch Target Identification (BTI). BTI marks valid indirect branch targets with landing pad instructions, causing the CPU to fault if an indirect branch lands elsewhere. This prevents attackers from redirecting indirect calls and jumps to arbitrary code.

### Platform

arm64 (requires ISA v8.5+)

### Toolchain

| Compiler | Minimal Version | Default Version | Flag |
|:---------|:----------------|:----------------|:-----|
| clang | 12.0 | - | `-mbranch-protection=bti` |
| gcc | 10.1 | - | `-mbranch-protection=bti` |


---

## ARM Memory Tagging Extension

- **Rule ID:** `arm-mte`
- **Implementation:** `ARMMTERule`

Checks for ARM Memory Tagging Extension (MTE). MTE assigns 4-bit tags to memory regions and pointers, detecting use-after-free and buffer overflows when tags mismatch during memory access.

### Platform

arm64 (requires ISA v8.5+)

### Toolchain

| Compiler | Minimal Version | Default Version | Flag |
|:---------|:----------------|:----------------|:-----|
| clang | 12.0 | - | `-march=armv8.5-a+memtag -fsanitize=memtag` |


---

## ARM Pointer Authentication

- **Rule ID:** `arm-pac`
- **Implementation:** `ARMPACRule`

Checks for ARM Pointer Authentication Code (PAC). PAC signs return addresses with a cryptographic key, detecting tampering when the signature is verified on function return. This prevents attackers from overwriting return addresses to hijack control flow.

### Platform

arm64 (requires ISA v8.3+)

### Toolchain

| Compiler | Minimal Version | Default Version | Flag |
|:---------|:----------------|:----------------|:-----|
| clang | 12.0 | - | `-mbranch-protection=pac-ret` |
| gcc | 10.1 | - | `-mbranch-protection=pac-ret` |


---

## Address Sanitizer

- **Rule ID:** `asan`
- **Implementation:** `ASANRule`

Checks for AddressSanitizer (ASan) instrumentation. ASan detects memory errors including buffer overflows, use-after-free, and memory leaks at runtime.

### Platform

arm64, x86, amd64, arm

### Toolchain

| Compiler | Minimal Version | Default Version | Flag |
|:---------|:----------------|:----------------|:-----|
| clang | 3.4 | - | `-fsanitize=address` |
| gcc | 5.1 | - | `-fsanitize=address` |


---

## ASLR Compatibility

- **Rule ID:** `aslr`
- **Implementation:** `ASLRRule`

Checks if the binary is compatible with Address Space Layout Randomization (ASLR). ASLR randomizes memory addresses at runtime, making it difficult for attackers to predict the location of code and data. This checks binary compatibility only, not system ASLR settings.

### Platform

arm, arm64, x86, amd64

### Toolchain

| Compiler | Minimal Version | Default Version | Flag |
|:---------|:----------------|:----------------|:-----|
| clang | 3.4 | 4.0 | `-fPIE -pie -z noexecstack` |
| gcc | 4.1 | 6.1 | `-fPIE -pie -z noexecstack` |
| rustc | 1.0 | 1.26 | `-C relocation-model=pie` |


---

## Control Flow Integrity

- **Rule ID:** `cfi`
- **Implementation:** `CFIRule`

Checks for Clang Control Flow Integrity (CFI) instrumentation. CFI validates that indirect calls and jumps target expected locations, preventing attackers from hijacking control flow through corrupted function pointers or vtables.

### Platform

x86, amd64, arm, arm64

### Toolchain

| Compiler | Minimal Version | Default Version | Flag |
|:---------|:----------------|:----------------|:-----|
| clang | 6.0 | - | `-fsanitize=cfi -flto -fvisibility=hidden` |


---

## FORTIFY_SOURCE

- **Rule ID:** `fortify-source`
- **Implementation:** `FortifySourceRule`

Checks for FORTIFY_SOURCE buffer overflow protection. This glibc feature replaces unsafe C library functions (strcpy, memcpy, sprintf, etc.) with bounds-checked variants at compile time.

### Platform

arm, arm64, x86, amd64

### Toolchain

| Compiler | Minimal Version | Default Version | Flag |
|:---------|:----------------|:----------------|:-----|
| clang | 12.0 | - | `-D_FORTIFY_SOURCE=3 -O1` |
| gcc | 12.1 | - | `-D_FORTIFY_SOURCE=3 -O1` |


---

## Full RELRO

- **Rule ID:** `full-relro`
- **Implementation:** `FullRELRORule`

Checks for full RELRO (Relocation Read-Only) protection. Full RELRO makes the Global Offset Table (GOT) read-only after initialization, preventing GOT overwrite attacks that redirect function calls to malicious code.

### Platform

x86, amd64, arm, arm64

### Toolchain

| Compiler | Minimal Version | Default Version | Flag |
|:---------|:----------------|:----------------|:-----|
| clang | 3.4 | 4.0 | `-Wl,-z,relro,-z,now` |
| gcc | 4.1 | 6.1 | `-Wl,-z,relro,-z,now` |
| rustc | 1.21 | 1.21 | `-C link-arg=-z -C link-arg=relro -C link-arg=-z -C link-arg=now` |


---

## Disallow dlopen

- **Rule ID:** `no-dlopen`
- **Implementation:** `NoDLOpenRule`

Checks if the shared library disallows being loaded via dlopen(). This prevents attackers from injecting the library into arbitrary processes.

### Platform

x86, amd64, arm, arm64

### Toolchain

| Compiler | Minimal Version | Default Version | Flag |
|:---------|:----------------|:----------------|:-----|
| clang | 3.4 | - | `-Wl,-z,nodlopen` |
| gcc | 4.1 | - | `-Wl,-z,nodlopen` |
| rustc | 1.0 | - | `-C link-arg=-z -C link-arg=nodlopen` |


---

## Core Dump Protection

- **Rule ID:** `no-dump`
- **Implementation:** `NoDumpRule`

Checks if the binary is excluded from core dumps. Disabling core dumps prevents sensitive data like cryptographic keys and passwords from being written to disk during crashes.

### Platform

arm64, x86, amd64, arm

### Toolchain

| Compiler | Minimal Version | Default Version | Flag |
|:---------|:----------------|:----------------|:-----|
| clang | 3.4 | - | `-Wl,-z,nodump` |
| gcc | 4.1 | - | `-Wl,-z,nodump` |


---

## Secure RPATH

- **Rule ID:** `no-insecure-rpath`
- **Implementation:** `NoInsecureRPATHRule`

Checks for insecure RPATH values that could enable library injection. RPATH takes precedence over system library paths, so relative paths or world-writable directories allow attackers to hijack library loading.

### Platform

arm, arm64, x86, amd64

### Toolchain

| Compiler | Minimal Version | Default Version | Flag |
|:---------|:----------------|:----------------|:-----|
| clang | 3.4 | - | `-Wl,-rpath,/absolute/path` |
| gcc | 4.1 | - | `-Wl,-rpath,/absolute/path` |
| rustc | 1.0 | - | `-C link-arg=-rpath -C link-arg=/absolute/path` |


---

## Secure RUNPATH

- **Rule ID:** `no-insecure-runpath`
- **Implementation:** `NoInsecureRUNPATHRule`

Checks for insecure RUNPATH values that could enable library injection. Relative paths, empty components, or world-writable directories in RUNPATH allow attackers to place malicious libraries that get loaded instead of legitimate ones.

### Platform

x86, amd64, arm, arm64

### Toolchain

| Compiler | Minimal Version | Default Version | Flag |
|:---------|:----------------|:----------------|:-----|
| clang | 3.4 | 4.0 | `-Wl,--enable-new-dtags -Wl,-rpath,/absolute/path` |
| gcc | 4.1 | 6.1 | `-Wl,--enable-new-dtags -Wl,-rpath,/absolute/path` |
| rustc | 1.74 | - | `-C link-arg=--enable-new-dtags -C link-arg=-rpath -C link-arg=/absolute/path` |


---

## Non-Executable Stack

- **Rule ID:** `nx-bit`
- **Implementation:** `NXBitRule`

Checks if the stack is marked non-executable (NX bit). This prevents stack-based buffer overflow exploits from executing shellcode placed on the stack.

### Platform

x86, amd64, arm, arm64

### Toolchain

| Compiler | Minimal Version | Default Version | Flag |
|:---------|:----------------|:----------------|:-----|
| clang | 1.0 | 1.0 | `-z noexecstack` |
| gcc | 3.4 | 3.4 | `-z noexecstack` |
| rustc | 1.0 | 1.0 | `-C link-arg=-Wl,-z,noexecstack` |


---

## Position Independent Executable

- **Rule ID:** `pie`
- **Implementation:** `PIERule`

Checks if the binary is compiled as a Position Independent Executable (PIE). PIE enables full ASLR by allowing the executable to be loaded at a random base address, making return-oriented programming (ROP) attacks significantly harder.

### Platform

arm, arm64, x86, amd64

### Toolchain

| Compiler | Minimal Version | Default Version | Flag |
|:---------|:----------------|:----------------|:-----|
| clang | 3.4 | 4.0 | `-fPIE -pie` |
| gcc | 4.1 | 6.1 | `-fPIE -pie` |
| rustc | 1.26 | 1.26 | `-C relocation-model=pie` |


---

## Partial RELRO

- **Rule ID:** `relro`
- **Implementation:** `RELRORule`

Checks for partial RELRO (Relocation Read-Only) protection. Partial RELRO reorders ELF sections to protect internal data structures and makes some segments read-only, but leaves the GOT writable for lazy binding.

### Platform

arm64, x86, amd64, arm

### Toolchain

| Compiler | Minimal Version | Default Version | Flag |
|:---------|:----------------|:----------------|:-----|
| clang | 3.4 | 3.9 | `-Wl,-z,relro` |
| gcc | 4.1 | 6.1 | `-Wl,-z,relro` |
| rustc | 1.21 | 1.21 | `-C link-arg=-z -C link-arg=relro` |


---

## SafeStack

- **Rule ID:** `safe-stack`
- **Implementation:** `SafeStackRule`

Checks for Clang SafeStack instrumentation. SafeStack separates the stack into a safe stack for return addresses and an unsafe stack for buffers, protecting control flow from stack buffer overflows.

### Platform

x86, amd64, arm, arm64

### Toolchain

| Compiler | Minimal Version | Default Version | Flag |
|:---------|:----------------|:----------------|:-----|
| clang | 3.7 | - | `-fsanitize=safe-stack` |


---

## Separate Code Segments

- **Rule ID:** `separate-code`
- **Implementation:** `SeparateCodeRule`

Checks if code and data are in separate memory pages. This prevents code pages from being writable and data pages from being executable, reducing the attack surface for memory corruption exploits.

### Platform

x86, amd64, arm, arm64

### Toolchain

| Compiler | Minimal Version | Default Version | Flag |
|:---------|:----------------|:----------------|:-----|
| clang | 6.0 | 6.0 | `-Wl,-z,separate-code` |
| gcc | 8.1 | 8.1 | `-Wl,-z,separate-code` |
| rustc | 1.31 | - | `-C link-arg=-z -C link-arg=separate-code` |


---

## Stack Canary Protection

- **Rule ID:** `stack-canary`
- **Implementation:** `StackCanaryRule`

Checks for stack canary (stack protector) instrumentation. Stack canaries detect buffer overflows by placing a guard value before the return address. If the canary is corrupted, the program terminates before exploitation can occur.

### Platform

x86, amd64, arm, arm64

### Toolchain

| Compiler | Minimal Version | Default Version | Flag |
|:---------|:----------------|:----------------|:-----|
| clang | 3.5 | 3.5 | `-fstack-protector-strong` |
| gcc | 4.9 | 4.9 | `-fstack-protector-strong` |


---

## Explicit Stack Size Limit

- **Rule ID:** `stack-limit`
- **Implementation:** `StackLimitRule`

Checks if an explicit stack size limit is set. Defining a maximum stack size helps prevent stack exhaustion attacks.

### Platform

x86, amd64, arm, arm64

### Toolchain

| Compiler | Minimal Version | Default Version | Flag |
|:---------|:----------------|:----------------|:-----|
| clang | 3.4 | - | `-Wl,-z,stack-size=<bytes>` |
| gcc | 4.1 | - | `-Wl,-z,stack-size=<bytes>` |


---

## Stripped Binary

- **Rule ID:** `stripped`
- **Implementation:** `StrippedRule`

Checks if the binary has been stripped of symbol tables and debug information. Stripping removes metadata that could help attackers understand the binary's structure and identify vulnerabilities.

### Platform

x86, amd64, arm, arm64

### Toolchain

| Compiler | Minimal Version | Default Version | Flag |
|:---------|:----------------|:----------------|:-----|
| clang | 1.0 | - | `-s` |
| gcc | 3.0 | - | `-s` |
| rustc | 1.59 | - | `-C strip=symbols` |


---

## Undefined Behavior Sanitizer

- **Rule ID:** `ubsan`
- **Implementation:** `UBSanRule`

Checks for Undefined Behavior Sanitizer (UBSan) instrumentation. UBSan detects undefined behavior such as integer overflows, null pointer dereferences, and misaligned accesses at runtime.

### Platform

arm64, x86, amd64, arm

### Toolchain

| Compiler | Minimal Version | Default Version | Flag |
|:---------|:----------------|:----------------|:-----|
| clang | 3.4 | - | `-fsanitize=undefined` |
| gcc | 5.1 | - | `-fsanitize=undefined` |


---

## x86 CET - Indirect Branch Tracking

- **Rule ID:** `x86-cet-ibt`
- **Implementation:** `X86CETIBTRule`

Checks for Intel Control-flow Enforcement Technology Indirect Branch Tracking (CET-IBT). IBT requires indirect branches to land on ENDBR instructions, preventing attackers from redirecting indirect calls and jumps to arbitrary code. Invalid branch targets trigger a control protection exception.

### Platform

x86, amd64

### Toolchain

| Compiler | Minimal Version | Default Version | Flag |
|:---------|:----------------|:----------------|:-----|
| clang | 10.0 | - | `-fcf-protection=full` |
| gcc | 8.1 | - | `-fcf-protection=full` |


---

## x86 CET - Shadow Stack

- **Rule ID:** `x86-cet-shstk`
- **Implementation:** `X86CETShadowStackRule`

Checks for Intel Control-flow Enforcement Technology Shadow Stack (CET-SS). Shadow Stack maintains a hardware-protected copy of return addresses, detecting ROP attacks when the shadow and regular stacks diverge on function return.

### Platform

x86, amd64

### Toolchain

| Compiler | Minimal Version | Default Version | Flag |
|:---------|:----------------|:----------------|:-----|
| clang | 10.0 | - | `-fcf-protection=full` |
| gcc | 8.1 | - | `-fcf-protection=full` |


---

## x86 Retpoline

- **Rule ID:** `x86-retpoline`
- **Implementation:** `X86RetpolineRule`

Checks for retpoline mitigation against Spectre v2 attacks. Retpoline replaces indirect branches with a return-based sequence that prevents speculative execution through the branch target buffer.

### Platform

x86, amd64

### Toolchain

| Compiler | Minimal Version | Default Version | Flag |
|:---------|:----------------|:----------------|:-----|
| clang | 6.0 | - | `-mretpoline` |
| gcc | 7.3 | - | `-mindirect-branch=thunk -mfunction-return=thunk` |

