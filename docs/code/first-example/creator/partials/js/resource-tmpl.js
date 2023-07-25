{{ $M := $.RESOURCE.Model }}
listTmpl = `
<table class="table">
  <thead>
    <tr>
      {{- range $F := $M.OrderedFields }}
      <th scope="col">{{ $F.Name }}</th>
      {{- end }}
    </tr>
  </thead>
  <tbody>
    {%#each . %}
    <tr>
      {{- range $F := $M.OrderedFields }}
      {{ if (eq $M.Index $F.Name) }}
      {{/* construct a table cell with a link to the element, crazy nested template systems herein be dragons */}}
      <td><a href="/{{ lower $.RESOURCE.Name }}?{{ lower $F.Name }}={% {{ $F.Name }} %}">{% {{ $F.Name }} %}</a></td>
      {{ else }}
      <td>{% {{ $F.Name }} %}</td>
      {{- end }}
      {{- end }}
    </tr>
    {%/each %}
  </tbody>
</table>
`

// add custom rendering
infoTmpl = `
`
