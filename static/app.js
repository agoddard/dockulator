;(function (calc) {
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
  }

}(window.calc = window.calc || {}));

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

