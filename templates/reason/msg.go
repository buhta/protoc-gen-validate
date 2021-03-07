package reason

const msgTpl = `
	{{- template "msgInner" . -}}
`

const msgInnerTpl = `
{{- range .NonOneOfFields }}
	// we have NonOneOfFields
	{{ renderConstants (context .) }}
{{ end }}
{{ range .OneOfs }}
	// we have OneOfs
	{{ template "oneOfConst" . }}
{{ end }}

{{ if disabled . }}
	// Validate is disabled for {{ simpleName . }}
	return;
{{- else -}}
	{{- range .NonOneOfFields}}
	let {{.Name}} = value => {	
		let errors = ref([]);

		{{ render (context .) }}

		List.length(errors^) == 0 ? Ok(value) : Error(errors); 
	};
{{ end -}}
{{ range .OneOfs }}
	{{ template "oneOf" . }}
{{- end -}}
{{- end }}
`
