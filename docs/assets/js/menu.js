window.toggleMenuItem = function (e) {
  e = e || window.event;
  var src = e.target || e.srcElement;

  if (src.nodeName === "SPAN") {
    src = src.children[0]
  }
  p1 = src.parentNode;
  p2 = p1.parentNode;

  children = p2.nextElementSibling;

  classes = src.className.split(" ");
  var open = classes.indexOf("fa-caret-down") >= 0
  if (open) {
    src.className = src.className.replace("down", "right")
    children.className += " d-none"
  } else {
    src.className = src.className.replace("right", "down")
    children.className = children.className.replace(" d-none", "")
  }
}
