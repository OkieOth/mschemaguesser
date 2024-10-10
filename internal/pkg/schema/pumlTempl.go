package schema

var pumlTemplateStr = `
@startuml
title Storage model for database: {{ .Database }}, collection: {{ .Collection }}\n\n
hide empty methods

footer Created with https://github.com/OkieOth/mschemaguesser

class "**{{ .MainType.Name }}**" as {{ .MainType.Name }} #FFFFFF {
  {{ range $index, $prop := .MainType.Properties -}}
  {{- $prop.AttribName }}: {{ $prop.ValueType }}
  {{- if $prop.IsArray }}[]{{ end }}
  {{- if not $prop.IsComplex }}<color:grey>    // {{ $prop.BsonType }}</color>{{ end }}
  {{ end -}}
}

{{if gt (len .MainType.Comments) 0}}
note top of {{ .MainType.Name }} 
    {{range .MainType.Comments}}
{{.}}
    {{end}}
end note
{{end}}

{{ $lastIndexOthers := LastIndexTypes .OtherComplexTypes -}}
{{- range $index, $type := .OtherComplexTypes -}}
  {{ if $type.IsDictionary -}}
class "**{{ $type.Name }}**" as {{ $type.Name }} <<Map>> #FFFFFF {
  valueType: {{ $type.DictValueType }}
} 

{{ $type.Name }} .. {{ $type.DictValueType }}

  {{ else }}
class "**{{ $type.Name }}**" as {{ $type.Name }} #FFFFFF {
  {{ range $index, $prop := $type.Properties -}}
  {{- $prop.AttribName }}: {{ $prop.ValueType }}
  {{- if $prop.IsArray }}[]{{ end }} 
  {{- if not $prop.IsComplex }}<color:grey>    // {{ $prop.BsonType }}</color>{{ end }}
  {{ end -}}
}  

  {{- end }}



{{ end }}

{{- range $index, $type := .Relations -}}
  {{ $type.Start }} *-- {{ $type.End }}
{{ end }}

@enduml
`
