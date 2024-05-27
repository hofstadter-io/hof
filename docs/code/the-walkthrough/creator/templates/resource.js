{{ template "js/resource-tmpl.js" . }}
{{ template "js/resource-utils.js"  . }}
{{ template "js/resource-lib.js"  . }}
{{ template "js/resource-init.js"  . }}

document.addEventListener("DOMContentLoaded", function(event) {
  run();
});

function run() {
  // fetch data
  fetchData().then(function(data) {
    // init html elements
    initElems(data)
  }).catch(function(error) {
    renderError(error)
  });
}
