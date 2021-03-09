package reason

const numConstTpl = `{{ $f := .Field }}{{ $r := .Rules -}}
{{- if $r.Const }}
	let {{ constantName . "Const" }} = {{ $r.GetConst }};
{{- end -}}
{{- if $r.Lt }}
	let {{ constantName . "Lt" }} = {{ $r.GetLt }};
{{- end -}}
{{- if $r.Lte }}
	let {{ constantName . "Lte" }} = {{ $r.GetLte }};
{{- end -}}
{{- if $r.Gt }}
	let {{ constantName . "Gt" }} = {{ $r.GetGt }};
{{- end -}}
{{- if $r.Gte }}
	let {{ constantName . "Gte" }} = {{ $r.GetGte }};
{{- end -}}
{{- if $r.In }}
	let {{ constantName . "In" }} = [
	{{- range $r.In -}}
		{{- sprintf "%v" . -}},
	{{- end -}}
	];
{{- end -}}
{{- if $r.NotIn }}
	let {{ constantName . "NotIn" }} = [
	{{- range $r.NotIn -}}
		{{- sprintf "%v" . -}},
	{{- end -}}
	];
{{- end -}}
`

const numTpl = `{{ $f := .Field }}{{ $r := .Rules -}}
{{- if $r.Const }}
		let errors = errors @ (value != {{ constantName . "Const" }} ? ["must be equal to " ++ string_of_int({{ constantName . "Gt" }})] : []);
{{- end -}}
{{- if $r.Gt }}
		let errors = errors @ (value > {{ constantName . "Gt" }} ? ["must be greater than " ++ string_of_int({{ constantName . "Gt" }})] : []);
{{- end -}}
{{- if $r.Lt }}
		let errors = errors @ (value < {{ constantName . "Lt" }} ? ["must be less than " ++ string_of_int({{ constantName . "Lt" }})] : []);
{{- end -}}
{{- if $r.Gte }}
		let errors = errors @ (value >= {{ constantName . "Gte" }} ? ["must be greater or equal to " ++ string_of_int({{ constantName . "Gte" }})] : []);
{{- end -}}
{{- if $r.Lte }}
		let errors = errors @ (value < {{ constantName . "Lte" }} ? ["must be less or equal to " ++ string_of_int({{ constantName . "Lte" }})] : []);
{{- end -}}
{{- if $r.NotIn }}
		// [TODO] print the array in the error message
		let errors = errors @ (List.exists( v => v == value, {{ constantName . "NotIn" }}) ? ["must be in the list: /*TODO: print array*/"] : []);
{{- end -}}
{{- if $r.In }}
		// [TODO] print the array in the error message
		let errors = errors @ (List.exists( v => v == value, {{ constantName . "In" }}) ? [] : ["must not be in the list: /*TODO: print array*/"]);
{{- end -}}
`
