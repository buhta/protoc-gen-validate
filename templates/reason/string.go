package reason

const stringConstTpl = `{{ $f := .Field }}{{ $r := .Rules -}}`

const stringTpl = `{{ $f := .Field }}{{ $r := .Rules -}}
{{- if $r.MinLen }}
		let errors = errors @ (String.length(value) >= {{ $r.GetMinLen }} ? [] : ["length must be min " ++ string_of_int({{ $r.GetMinLen }})]);
{{- end -}}
{{- if $r.MaxLen }}
		let errors = errors @ (String.length(value) <= {{ $r.GetMaxLen }} ? [] : ["length must be max " ++ string_of_int({{ $r.GetMaxLen }})]);
{{- end -}}
`
