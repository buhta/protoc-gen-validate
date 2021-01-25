package ocaml

const requiredTpl = `{{ $f := .Field }}
{{- if .Rules.GetRequired }}
      let {{$f.Descriptor.Name}} value = match value with | Some value -> value | None  -> Error("field is required");
{{- end -}}
`
