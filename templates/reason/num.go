package reason

const numConstTpl = `{{ $f := .Field }}{{ $r := .Rules -}}
{{- if $r.Const }}
	let {{ constantName . "Const" }} = {{ $r.GetConst }};
{{- end -}}
{{- if $r.Lt }}
	let {{ constantName . "Lt" }} = {{ $r.GetLt }};
{{- end -}}
{{- if $r.Gt }}
	let {{ constantName . "Gt" }} = {{ $r.GetGt }};
{{- end -}}`

const numTpl = `{{ $f := .Field }}{{ $r := .Rules -}}
{{- if $r.Gt }}
		if (value > {{ constantName . "Gt" }}) {
			errors := errors^ @ ["must be greater than " ++ string_of_int({{ constantName . "Gt" }})];
		};
{{- end -}}
{{- if $r.Lt }}
		if (value < {{ constantName . "Lt" }}) {
			errors := errors^ @ ["must be Lt than " ++ string_of_int({{ constantName . "Lt" }})];
		};
{{- end -}}
`
