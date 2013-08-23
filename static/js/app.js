;(function (websock) {
  if ("WebSocket" in window) {
    ws = new WebSocket("ws://dockulator.com/websock");
    ws.onopen = function () {};
    ws.onclose = function () {
      // This is when the server closes the connection
    };
    ws.onerror = function () {};

    ws.sendmsg = function (message) {
      ws.send(message);
    };

    ws.onmessage = function (event) {
      if (event.data == 'ping') { return; }
      var obj = JSON.parse(event['data']);
      switch (obj['type']) {
        case 'new':
          websock.newCalculation(obj['data']);
          break;
        case 'update':
          websock.updateCalculation(obj['data']);
          break;
        case 'error':
          error(obj['data']);
          break;
        default:
          break;
      }
    };

    websock.ws = ws;
  } else {
    websock.ws = {};
  }

  websock.updateCalculation = function (data) {
    var c = new calc.Calculation(data);
    if (c.el().length == 0) {
      c.render()
    } else {
      c.el().remove();
      c.render();
    }
  };

  websock.addCalculation = function (data) {
    var c = new calc.Calculation(data);
    c.render();
  };

  websock.newCalculation = function (data) {
    data.answer = 'Calculating...';
    var c = new calc.Calculation(data);
    c.render('warning');
  };

  websock.error = function (error) {
  };

}(window.websock = window.websock || {}));

;(function (calc, $) {
  calc.Calculation = function (obj) {
    this.id = 'id-' + obj['id'];
    this.calculation = obj['calculation'];
    this.language = obj['language'];
    this.displayLang = getDisplayLanguage(obj['language']);
    this.answer = obj['answer'];
    this.os = obj['os'];
  };

  calc.Calculation.prototype.render = function (className) {
    var c = document.createElement('td'),
        ctext = document.createTextNode(this.calculation),
        a = document.createElement('td'),
        atext = document.createTextNode(this.answer),
        o = document.createElement('td'),
        otext = document.createTextNode(this.os),
        l = document.createElement('td'),
        ltext = document.createTextNode(this.displayLang),

        tableEl = document.getElementById("calcbody"),
        trEl = document.createElement('tr');

    c.appendChild(ctext);
    a.appendChild(atext);
    o.appendChild(otext);
    l.appendChild(ltext);

    // Do something cute for the class names
    l.className = this.language;

    trEl.appendChild(c);
    trEl.appendChild(a);
    trEl.appendChild(o);
    trEl.appendChild(l);
    trEl.className = className;
    trEl.id = this.id;

    tableEl.insertBefore(trEl, tableEl.firstChild);
  };

  calc.Calculation.el = function() { return $('#' + this.id); };

  var getDisplayLanguage = function (language) {
    switch (language) {
      case 'rb':
        return 'Ruby';
      default:
        return 'Unknown language';
    }
  };

}(window.calc = window.calc || {}, jQuery));


// Init stuff
;(function ($) {
  var hideSel = '#hideme',
      calcInputSel = '#input-calculation',
      formSel = '#new-calculation',
      hideKey = 'hideInfo';

  var hasLocalStorage = function () {
    try {
      return 'localStorage' in window && window['localStorage'] !== null;
    } catch(e){
      return false;
    }
  };

  var getCalculation = function () {
    input = $(calcInputSel).val();
    return input;
  };

  var alwaysHide = function () {
    if (hasLocalStorage()) {
      localStorage.setItem(hideKey, true);
    }
  };

  var alwaysShow = function () {
    if (hasLocalStorage()) {
      localStorage.removeItem(hideKey);
    }
  };

  $(hideSel).on('click', function () { 
    if (this.innerHTML === "Hide info") {
      alwaysHide();
      this.innerHTML = "Show info";
    } else {
      alwaysShow();
      this.innerHTML = "Hide info";
    }
  });

  if (hasLocalStorage()) {
    if (localStorage.getItem(hideKey)) {
      $(hideSel).click();
    }
  }

  $(formSel).on('submit', function (event) {
    $.ajax({
      type: 'POST',
      url: '/calculations',
      data: {
        'calculation': getCalculation()
      },
      beforeSend: function () {
        $(calcInputSel).val('');
      }
    }).done(function (data) {});
    return false;
  });

}(jQuery));
