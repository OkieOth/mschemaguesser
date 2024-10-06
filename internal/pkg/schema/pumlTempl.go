package schema

var pumlTemplateStr = `
@startuml
title Storage model for database: {{ .Database }}, collection: {{ .Collection }}\n\n

footer Created with https://github.com/OkieOth/mschemaguesser

{{if gt (len .MainType.Comments) 0}}

class {{ .MainType.Name }} {

}

note top of {{ . MainType.Name }} 
    {{range .MainType.Comments}}
{{.}}
    {{end}}
end note
{{end}}

{{ $lastIndexOthers := LastIndexTypes .OtherComplexTypes -}}
{{- range $index, $type := .OtherComplexTypes -}}
class {{ $type.Name }} {
}{{ if ne $index $lastIndexOthers }},{{ end }}
{{ end }}

{{- range $index, $type := .Relations -}}
  $type.Start o- $type.End
{{ end }}

@enduml
`
