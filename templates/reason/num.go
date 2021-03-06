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
	let {{ $f.Name }} = value => value > {{ constantName . "Gt" }} ? Ok : Error("must be greater than " ++ {{ constantName . "Gt" }});
{{- end -}}
`
