package schema

var schemaTemplateStr = `
{
  "$schema": "http://json-schema.org/draft-07/schema#",
  "title": "{{ .MainType.Name }}",
  "description": "Storage model for database: {{ .Database }}, collection: {{ .Collection }}",
  "version": "0.0.0",
  "x-model-type": "mongodb-storage-model",
  {{ if .MainType.Count.IsSet -}}
  "x-collection-elem-count": "{{ .MainType.Count.Value }}",{{- end -}}
  {{if gt (len .MainType.Comments) 0}}
    {{range .MainType.Comments}}
  "x-processing-comment": "{{.}}",
    {{end}}
  {{end}}
  "type": "object",
  "properties": {
    {{ $lastIndexProps := LastIndexProps .MainType.Properties -}}
    {{- range $index, $prop := .MainType.Properties -}}
    "{{- $prop.AttribName }}": { {{ if $prop.IsArray }}
      "type": "array",
      "items": {
        "x-bson-type": "{{ $prop.BsonType }}",
        {{ if $prop.IsComplex -}} "$ref": "#/definitions/{{ $prop.ValueType }}"
        {{ else -}} "type": "{{ $prop.ValueType }}"
        {{- end }}
      }
      {{ else }}
      {{ if ne $prop.Comment "" -}}
      "x-comment": "{{ $prop.Comment }}",{{- end }}
      "x-bson-type": "{{ $prop.BsonType }}",
      {{ if ne $prop.Format "" -}}  "format": "{{ $prop.Format }}",
      {{- end }}
      {{ if $prop.IsComplex -}} "$ref": "#/definitions/{{ $prop.ValueType }}"
      {{- else -}} "type": "{{ $prop.ValueType }}"
      {{- end }}
      {{- end }}
    }{{ if ne $index $lastIndexProps }},{{ end }}
    {{ end }}
  },

  "definitions": {
    {{ $lastIndexOthers := LastIndexTypes .OtherComplexTypes -}}
    {{- range $index, $type := .OtherComplexTypes -}}
    "{{ $type.Name }}": {
      "type": "object",
      "x-dict": {{ $type.IsDictionary }},
      {{ if $type.IsDictionary -}}
        "additionalProperties": {
            "$ref": "#/definitions/{{ $type.DictValueType }}"
      }
      {{ else }}
      "properties": {
        {{- $lastIndexProps := LastIndexProps $type.Properties -}}
        {{- range $index, $prop := $type.Properties }}
        "{{ $prop.AttribName }}": { {{ if $prop.IsArray -}}
          "type": "array",
          "items": {
            "x-bson-type": "{{ $prop.BsonType }}", {{ if $prop.IsComplex -}}
            "$ref": "#/definitions/{{ $prop.ValueType }}"
            {{- else -}}
            "type": "{{ $prop.ValueType }}"
            {{- end }}
          }
          {{- else -}}
          {{ if ne $prop.Comment "" -}}
          "x-comment": "{{ $prop.Comment }}",
          {{- end }}
          "x-bson-type": "{{ $prop.BsonType }}",
          {{ if ne $prop.Format "" -}}  "format": "{{ $prop.Format }}",
          {{- end }}
          {{- if $prop.IsComplex -}} 
          "$ref": "#/definitions/{{ $prop.ValueType }}"
          {{- else -}} 
          "type": "{{ $prop.ValueType }}"
          {{- end }}
          {{- end }}
        }{{ if ne $index $lastIndexProps }},{{ end -}}
        {{- end }}
      }
      {{- end }}
    }{{ if ne $index $lastIndexOthers }},{{ end }}
    {{ end }}
  }
}
`
