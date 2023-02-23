function renderData(data, elem) {
  var solo = isSolo(data);

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

function initElems(data) {
  var action = determineAction();
  console.log(action , data)

  // init data display
  elem = document.getElementById("data");
  renderData(data, elem);

  // init mutation button
  el = document.getElementById("modal-btn");
  el.innerHTML = action;

  // populate the form
  if (action === "update") {
    console.log("setting form")
    el = getForm();
    {{ range $.RESOURCE.Model.OrderedFields }}
    el.{{ .Name }}.value = data.{{ .Name }};
    {{- end }}
  }
}
