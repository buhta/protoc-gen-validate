package reason

const fileTpl = `// Code generated by protoc-gen-validate. DO NOT EDIT.
// source: {{ .File.InputPath }}

{{ if isOfFileType . }}
module Validators = {
{{ range .AllMessages -}}
	{{- template "msg" . -}}
{{- end }}
}
{{ end }}
`
