package ocaml

const msgTpl = `
{{ if not (ignored .) -}}
	{{- template "msgInner" . -}}
{{- end -}}
`

const msgInnerTpl = `
	{{- range .NonOneOfFields }}
		{{ renderConstants (context .) }}
	{{ end }}
	{{ range .OneOfs }}
		{{ template "oneOfConst" . }}
	{{ end }}

	{{ if disabled . }}
		// Validate is disabled for {{ simpleName . }}
		return;
	{{- else -}}
	{{ range .NonOneOfFields -}}
		{{ render (context .) }}
	{{ end -}}
	{{- end }}
	`

// {{ range .OneOfs }}
// 	{{ template "oneOf" . }}
// {{- end -}}
