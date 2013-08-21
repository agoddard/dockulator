;(function (websock) {
  websock.ws = new WebSocket("ws://dockulator.com/websock");

  var updateCalculation = function (data) {
    var el = $('#id-' + data['id']);
    if (el.length == 0) {
      addCalculation(data);
      return
    }
    el.parentNode.removeChild(el);
    addCalculation(data);
  };
  var addCalculations = function (data) {
    calc.insert.apply(null, data);
  };
  var addCalculation = function (data) {
    calc.insert.apply(null, [data]);
  };
  var error = function (error) {
    console.log("Got an error from the server:", error);
  };

  websock.ws.onmessage = function (event) {
    if (event.data == 'ping') {
      return;
    }
    console.log(event.data);
    var obj = JSON.parse(event['data']);
    switch (obj['type']) {
      case 'new':
        console.log("new");
        addCalculation(obj['data']);
        break;
      case 'update':
        console.log("update");
        updateCalculation(obj['data']);
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
        frag = document.createDocumentFragment(),
        obj, el, tmpl;

    for (;i < arguments.length; i++) {
      obj = arguments[i];
      el = document.createElement('tr');
      el.className = 'calculation';
      el.id = obj['id'];
      console.log(el);
      tmpl = calc.render(obj);
      el.appendChild(tmpl);
      frag.appendChild(el);
    }
    tableEl.insertBefore(frag, tableEl.firstChild);
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

