{{- range $i, $r := .Rules -}}
{{- if $i}}

---

{{end -}}
## {{$r.Name}}

- **Rule ID:** `{{$r.ID}}`
- **Implementation:** `{{$r.StructName}}`

{{$r.Description}}

### Platform

{{$r.Platform}}

### Toolchain

{{if $r.Compilers -}}
| Compiler | Minimal Version | Default Version | Flag |
|:---------|:----------------|:----------------|:-----|
{{range $r.Compilers -}}
| {{.Name}} | {{.MinVersion}} | {{.DefaultVersion}} | `{{.Flag}}` |
{{end -}}
{{- else -}}
No specific compiler requirements.
{{end -}}
{{- end}}
