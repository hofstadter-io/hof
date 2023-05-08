---
title: "Variables & Scopes"
weight: 21
---

{{<lead>}}
The underlying Go text/template system has semantics which
will impact how variables are scoped and can be accessed.
This page will help you understand and work through these.
{{</lead>}}



## Range and With statements

A range or with statement will hide any `.Value` paths
because it creates a new scope.
To access variables outside, assign them to `$Value`
like in this example.

We often recommend that you capture your important
variables at the start of a template or partial
as a best practice you should adopt.

```
{{ $DM := .DM }}
{{ $M := .Model }}
CREATE TABLE {{ snake $M.Plural }} (
{{ range $K, $F := $M.Fields }}
{{ if ne $K "$hof" }}
{{ template "sql/create-field.sql" (dict "Field" $F "Model" $M "DM" $DM) }}
{{ end }}
{{ end }}
);
```

## Passing variables to sub-templates (partials)

In the same way, a call to a partial template will create a new scope.
Your top-level or saved variables will not be accessible.
These partial templates only accept one value,
so to pass multiple, we use the `dict` helper.
You can see this in the above example.

The partial template referenced:

```
{{ snake .Field.Name }} {{ if .Field.sqlType -}}
{{ .Field.sqlType }}{{ else }}{{ .Field.Type -}}{{end}}{{ with .Field.Default }} DEFAULT {{.}}{{ end }},
{{ if .Field.SQL.PrimaryKey }}constraint {{ .Model.Plural }}_pkey primary key ({{ snake .Field.Name }}),{{ end }}
```
