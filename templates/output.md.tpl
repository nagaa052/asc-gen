# {{ .Title }}

{{ range $label, $lines := .Content }}
## {{ $label }}
|  Name  |  Description  | UpdatedAt |
| :--- | :--- | :--- |
{{ range $i, $line := $lines -}}
| [{{ $line.Name }}]({{ $line.URL }}) | {{ $line.Description }} | {{ $line.UpdatedAt }} |
{{ end -}}
{{ end }}