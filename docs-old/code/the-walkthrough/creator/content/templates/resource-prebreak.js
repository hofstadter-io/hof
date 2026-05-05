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

  elem.innerHTML = JSON.stringify(data);

}
