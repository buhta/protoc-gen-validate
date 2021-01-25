package ocaml

const numTpl = `{{ $f := .Field }}{{ $r := .Rules -}}
{{- if $r.Const }}
      - IdenticalTo:
          value: {{ $r.GetConst }}{{ ocamlTypeLiteralSuffixFor . }}
{{- end -}}
{{- if and (or $r.Lt $r.Lte) (or $r.Gt $r.Gte)}}
      - TODONumLteGteRange: ~
{{- else -}}
{{- if $r.Lt }}
      - LessThan: {{ $r.GetLt }}{{ ocamlTypeLiteralSuffixFor . }}
{{- end -}}
{{- if $r.Lte }}
      - LessThanOrEqual: {{ $r.GetLte }}{{ ocamlTypeLiteralSuffixFor . }}
{{- end -}}
{{- if $r.Gt }}
      let {{$f.Descriptor.Name}} value =  value > if {{ $r.GetGt }} then Ok(value) else Error("value should be greater than {{ $r.GetGt }}");
{{- end -}}
{{- if $r.Gte }}
      - GreaterThanOrEqual: {{ $r.GetGte }}{{ ocamlTypeLiteralSuffixFor . }}
{{- end -}}
{{- end -}}
`
