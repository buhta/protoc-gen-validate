package ocaml

const stringTpl = `{{ $f := .Field }}{{ $r := .Rules -}}
{{- if $r.Const }}
      - EqualTo: {{ $r.GetConst }}
{{- end -}}
{{- if $r.Len }}
      - Length: {{ $r.GetLen }}
        min: {{ $r.GetLen }}
        max: {{ $r.GetLen }}
{{- end -}}
{{- if $r.MinLen }}
	let {{$f.Descriptor.Name}} value =
		if String.length value > {{$r.GetMinLen}} then Ok(value) : Error("length of {{$f.Descriptor.Name}} should be greater than {{ $r.GetMinLen }}");
{{- end -}}
{{- if $r.LenBytes }}
      - Length: {{ $r.GetLenBytes }}
        min: {{ $r.GetLenBytes }}
        max: {{ $r.GetLenBytes }}
{{- end -}}
{{- if or $r.MinBytes $r.MaxBytes }}
      - Length:
        {{- if $r.MinBytes }}
        min: {{ $r.GetMinBytes }}
        {{- end -}}
        {{- if $r.MaxBytes }}
        max: {{ $r.GetMaxBytes }}
        {{- end -}}
{{- end -}}
{{- if $r.Pattern }}
      - Regex:
        pattern: "/{{ ocamlStringEscape $r.GetPattern }}/"
{{- end -}}
{{- if $r.Prefix }}
      - TODOStringPrefix: {{ ocamlStringEscape $r.GetPrefix }}
{{- end -}}
{{- if $r.Contains }}
      - TODOStringContains: {{ ocamlStringEscape $r.GetContains }}
{{- end -}}
{{- if $r.NotContains }}
      - TODOStringNotContains: {{ ocamlStringEscape $r.GetNotContains }}
{{- end -}}
{{- if $r.Suffix }}
      - TODOStringSuffix: {{ ocamlStringEscape $r.GetSuffix }}
{{- end -}}
{{- if $r.GetEmail }}
      - Email: ~
{{- end -}}
{{- if $r.GetAddress }}
      - Hostname: ~
{{- end -}}
{{- if $r.GetHostname }}
      - Hostname: ~
{{- end -}}
{{- if $r.GetIp }}
      - Ip:
        version: all
{{- end -}}
{{- if $r.GetIpv4 }}
      - Ip:
        version: 4
{{- end -}}
{{- if $r.GetIpv6 }}
      - Ip:
        version: 6
{{- end -}}
{{- if $r.GetUri }}
      - Url: ~
{{- end -}}
{{- if $r.GetUriRef }}
      - Url: ~
        relativeProtocol: true
{{- end -}}
{{- if $r.GetUuid }}
      - Uuid: ~
{{- end -}}
`
