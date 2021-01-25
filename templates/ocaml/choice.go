package ocaml

const choiceTpl = `{{ $f := .Field }}{{ $r := .Rules -}}
{{- if $r.In }}
      - Choice:
          choices:
            {{- range $f.In }}
            - {{ sprintf "%v" . }}{{ ocamlTypeLiteralSuffixFor $ }}
            {{- end }}
		  # message:  .
{{- end -}}
{{- if $r.NotIn }}
      - NotInChoice:
          choices:
            {{- range $f.NotIn }}
            - {{ sprintf "%v" . }}{{ ocamlTypeLiteralSuffixFor $ }}
            {{- end }}
{{- end -}}
`
