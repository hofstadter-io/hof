function addHeaderLinks() {
  var content = document.getElementById("main-content");

  var h2s = content.querySelectorAll("h2");
  var h3s = content.querySelectorAll("h3");

  elems = [...h2s, ...h3s];

  elems.forEach(elem => {
    if (elem.parentNode == content) {
      var a = document.createElement("a");
      a.className = "ms-2 fas fa-link fa-xs";
      a.href = window.location.pathname + "#" + elem.id
      elem.className = "anchor"
      elem.appendChild(a);
    }
  })
}

addHeaderLinks();
