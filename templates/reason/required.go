package reason

const requiredTpl = `{{ $f := .Field }}
{{- if .Rules.GetRequired }}
		let errors = errors @ switch value {
				| None => ["required field"]
				| _ => []
			}
{{- end -}}`
