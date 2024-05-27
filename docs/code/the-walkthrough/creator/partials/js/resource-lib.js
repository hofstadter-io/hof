{{ $R := $.RESOURCE }}
function delete{{ $R.Name }}({{ $R.Model.Index }}) {
  url = makeUrl()
  var req = new Request(url, { method: "DELETE" })
  return fetch(req).then(function(response) {
    if (response.status >= 400) {
      throw data;
    } else {
      window.location.href = "/{{ lower $R.Name }}"
    }
  });
}

function create{{ $R.Name }}({{ template "js/resource-lib-args.js" . }}) {
  console.log("create {{ $R.Name }}")
  url = "/api" + window.location.pathname
  var req = new Request(url, { method: "POST" })
  return fetch(req).then(function(response) {
    if (response.status >= 400) {
      throw data;
    } else {
      window.location.href = "/{{ lower $R.Name }}?{{ lower $R.Model.Index }}=" + {{ $R.Model.Index }}
    }
  });
}

var myHeaders = new Headers();
myHeaders.append('Content-Type', 'application/json');

function update{{ $R.Name }}(
  {{- range $i, $F := $R.Model.OrderedFields -}}
  {{- if (gt $i 0) }}, {{ end }}{{ $F.Name -}}
  {{- end -}}
) {
  console.log("update {{ $R.Name }}")
  url = "/api" + window.location.pathname + "?{{$R.Model.Index}}=" + {{$R.Model.Index}}
  var req = new Request(url, {
    method: "PATCH",
    headers: myHeaders,
    body: JSON.stringify({
      {{- range $i, $F := $R.Model.OrderedFields -}}
      {{- if (gt $i 0) }}, {{ end }}{{ $F.Name -}}
      {{- end -}}
    })
  })
  return fetch(req).then(function(response) {
    if (response.status >= 400) {
      throw data;
    } else {
      window.location.href = "/{{ lower $R.Name }}"
    }
  });
}

function {{ camel $R.Name }}Submit() {
  var action = determineAction()
  var form = getForm()
  if (action === "create") {
    create{{ $R.Name }}({{ template "js/resource-lib-pass.js" . }})
  }
  if (action === "update") {
    update{{ $R.Name }}({{ template "js/resource-lib-pass.js" . }})
  }
}

// determine what we are doing, based on existance of query params
function determineAction() {
  if (window.location.search === "") {
    return "create"
  } else {
    return "update"
  }
}

function getForm() {
  return {
    {{ range $.RESOURCE.Model.OrderedFields }}
    {{ .Name }}: document.getElementById("{{ .Name }}Input"),
    {{- end }}
    form: document.getElementById("modal-form")
  }
}

function isSolo(data) {
  return !(data.length && data.length >= 0);
}
