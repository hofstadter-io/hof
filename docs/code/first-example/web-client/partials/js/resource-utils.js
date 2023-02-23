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

function fetchData() {
  url = makeUrl();
  // call API
  return fetch(url).then(function(response) {
    return response.json().then(function(data){
      if (response.status >= 400) {
        throw data;
      }
      return data;
    })
  })
}

function renderError(error, elem) {
  elem = document.getElementById("data");
  elem.innerHTML = error.message
}

