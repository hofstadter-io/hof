These are extra instruction files associated with the project and files you have available.
IMPORTANT: these are highly relevant by nature of contextual relevance and curated authorship.

<extra-instructions>
{{ range $path,$content := .agentsMd }}
<file path="{{$path}}">
{{$content}}
</file>
{{ end}}
</extra-instructions>