;(function (calc) {
  calc.render = function (obj) {
    var c = document.createElement('td'),
        ctext = document.createTextNode(obj['calculation']),
        a = document.createElement('td'),
        atext = document.createTextNode(obj['answer']),
        o = document.createElement('td'),
        otext = document.createTextNode(obj['os']),
        l = document.createElement('td'),
        ltext = document.createTextNode(obj['language']),
        frag = document.createDocumentFragment();

    c.appendChild(ctext);
    a.appendChild(atext);
    o.appendChild(otext);
    l.appendChild(ltext);

    frag.appendChild(c);
    frag.appendChild(a);
    frag.appendChild(o);
    frag.appendChild(l);
    return frag
  };

  calc.insert = function() {
    var i = 0, 
        len = arguments.length,
        tableEl = document.getElementById("calcbody"),
        frag = document.createDocumentFragment();

    for (;i < arguments.length; i++) {
      var el = document.createElement('tr');
      el.className = "calculation"
      var tmpl = calc.render(arguments[i]);
      el.appendChild(tmpl);
      frag.appendChild(el);
    }
    tableEl.appendChild(frag);
  }

}(window.calc = window.calc || {}));

;(function () {
  var ws = new WebSocket("ws://localhost:5000/websock");

  var cleanInput = function (data) {
    return data.substr(1, data.length-2);
  }

  var getJson = function (data) {
    return JSON.parse(atob(cleanInput(data)));
  }

  ws.onmessage = function (event) {
    var data = getJson(event.data);
    calc.insert.apply(null, data);
  };

}());
