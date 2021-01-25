package ocaml

const enumTpl = `{{ $f := .Field }}{{ $r := .Rules -}}
{{- if $r.Const }}
      - IdenticalTo:
          value: {{ $r.Const }}
{{- end -}}
{{- if $r.GetDefinedOnly }}
  let {{$f.Descriptor.Name}} value = Ok(value);
{{- end -}}
{{- if $r.In }}
      - Choice: # Enum.In
          choices:
            {{- range $r.In }}
            - {{ sprintf "%v" . }}
            {{- end }}
          # message:  .
{{- end -}}
{{- if $r.NotIn }}
      - NotInChoice: # Enum.NotIn
          choices:
            {{- range $r.NotIn }}
            - {{ sprintf "%v" . }}
            {{- end }}
          # message:  .
{{- end -}}
`
