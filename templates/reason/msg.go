package reason

const msgTpl = `
	{{- template "msgInner" . -}}
`

const msgInnerTpl = `
{{- range .NonOneOfFields }}
	{{- renderConstants (context .) -}}
{{- end -}}
{{ range .OneOfs }}
	{{- template "oneOfConst" . -}}
{{ end }}

{{ if disabled . }}
	// Validate is disabled for {{ simpleName . }}
	return;
{{- else -}}
	{{- range .NonOneOfFields}}
	// we have NonOneOfField {{.}}
	let {{.Name}} = value => {	
		let errors = ref([]);
		{{- render (context .)}}
		
		List.length(errors^) == 0 ? Ok(value) : Error(errors^); 
	};
{{ end -}}
{{ range .OneOfs }}
	// we have OneOfs {{.}}
	{{ template "oneOf" . }}
{{- end -}}
{{- end }}
`
