package reason

const enumConstTpl = `{{ $ctx := . }}{{ $f := .Field }}{{ $r := .Rules -}}
{{- if $r.In }}
	let {{ constantName . "In" }} = [
		{{- range $r.In }}
		{{.}},
		{{- end }}
	];
{{- end -}}
{{- if $r.NotIn }}
	let {{ constantName . "NotIn" }} = [
		{{- range $r.NotIn }}
		{{.}},
		{{- end }}
	];
{{- end -}}`

const enumTpl = `{{ $f := .Field }}{{ $r := .Rules -}}
{{- if $r.Const }}
		let errors = errors @ (value != {{ $r.GetConst }} ? ["must be equal to " ++ string_of_int({{ $r.GetConst }})] : []);
{{- end -}}
{{- if $r.In }}
		let errors = errors @ (List.exists( v => v == value, {{ constantName . "In" }}) ? [] : ["must be in the list: /*TODO: print array*/"]);
{{- end -}}
{{- if $r.NotIn }}
		let errors = errors @ (List.exists( v => v == value, {{ constantName . "NotIn" }}) ? ["must not be in the list: /*TODO: print array*/"] : []);
{{- end -}}
`
