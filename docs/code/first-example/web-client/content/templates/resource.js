document.addEventListener("DOMContentLoaded", function(event) {
  makeReq();
});

function makeUrl() {
  loc = window.location
  const urlSearchParams = new URLSearchParams(loc.search);
  const params = Object.fromEntries(urlSearchParams.entries());

  id = ""
  if (loc.search !== "") {
    id = "/" + params.{{ lower $.RESOURCE.Model.Index }}
  }

  url = "/api" + loc.pathname + id
  return url
}

function makeReq() {
  // where we put content
  elem = document.getElementById("data");

  url = makeUrl();

  // call API
  fetch(url).then(function(response) {
    return response.json().then(function(data){
      if (response.status >= 400) {
        throw data;
      }
      return data;
    })
  }).then(function(data) {
    renderData(data, elem)
  }).catch(function(error) {
    renderError(error, elem)
  });
}

function renderError(error, elem) {
  elem.innerHTML = error.message
}

function renderData(data, elem) {
  console.log("data", data)
  var solo = !(data.length && data.length >= 0);
  console.log("han: ", solo)

  var tmpl = infoTmpl;
  if (!solo) {
    tmpl = listTmpl;
  }
  t = tmpl.replaceAll("{%", "{" + "{").replaceAll("%}", "}" + "}");
  var template = Handlebars.compile(t);
  elem.innerHTML = template(data);

  if (solo) {
    btn = document.getElementById("del-btn")
    btn.onclick = delete{{ $.RESOURCE.Name }}
  }
}

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

function delete{{ $.RESOURCE.Name }}({{ $.RESOURCE.Model.Index }}) {
  url = makeUrl()
  console.log("Del: ", url)
  var req = new Request(url, { method: "DELETE" })
  return fetch(req).then(function(response) {
    if (response.status >= 400) {
      throw data;
    } else {
      window.location.href = "/{{ lower $.RESOURCE.Name }}"
    }
  });
}
