package reason

const messageTpl = `{{ $f := .Field }}{{ $r := .Rules }}
{{- if .MessageRules.GetSkip }}
    	// Skipping validation for {{ $f.Name }}
{{- else }}
		{{- template "required" . }}
{{- end -}}
`
