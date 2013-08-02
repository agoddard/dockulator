;(function (websock) {
  websock.ws = new WebSocket("ws://localhost:5000/websock");

  var initialData = function (data) {
      calc.insert.apply(null, data);
  }
  var error = function (error) {
    console.log("Got an error from the server:", error);
  }

  websock.ws.onmessage = function (event) {
    if (event.data == 'ping') {
      return;
    }
    var obj = JSON.parse(event['data']);
    switch (obj['type']) {
      case 'initialData':
        initialData(obj['data']);
        break;
      case 'error':
        error(obj['data']);
        break;
      default:
        console.log("Got a message, unsure how to proceed")
        console.log(event);
        break;
    }
  };

  websock.ws.onopen = function () {};

  websock.ws.onclose = function () {
    // This is when the server closes the connection
  };

  websock.ws.onerror = function () {};

  websock.ws.sendmsg = function (message) {
    websock.ws.send(message);
  };

}(window.websock = window.websock || {}));

;(function (calc, $) {
  calc.getDisplayLanguage = function (language) {
    switch (language) {
      case 'rb':
        return 'Ruby';
      default:
        return 'Unknown language';
    }
  };

  calc.render = function (obj) {
    var c = document.createElement('td'),
        ctext = document.createTextNode(obj['calculation']),
        a = document.createElement('td'),
        atext = document.createTextNode(obj['answer']),
        o = document.createElement('td'),
        otext = document.createTextNode(obj['os']),
        l = document.createElement('td'),
        ltext = document.createTextNode(calc.getDisplayLanguage(obj['language'])),
        frag = document.createDocumentFragment();

    c.appendChild(ctext);
    a.appendChild(atext);
    o.appendChild(otext);
    l.appendChild(ltext);

    frag.appendChild(c);
    frag.appendChild(a);
    frag.appendChild(o);
    frag.appendChild(l);
    return frag;
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
  };

  calc.getCalculation = function () {
    var calcInputId = '#input-calculation';
    input = $(calcInputId).val();
    return input;
  };

  var calcFormId = '#new-calculation';

  $('#new-calculation').on('submit', function (event) {
    $.ajax({
      type: 'POST',
      url: '/calculations',
      data: {
        'calculation': calc.getCalculation()
      }
    }).done(function (data) {
      console.log(data);
    });
    return false;
  });

}(window.calc = window.calc || {}, jQuery));

