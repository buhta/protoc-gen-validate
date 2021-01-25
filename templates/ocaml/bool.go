package ocaml

const boolTpl = `{{ $f := .Field }}{{ $r := .Rules -}}
{{- if $r.Const }}
      - IdenticalTo:
           value: {{ $r.GetConst }}
{{- end -}}
`
