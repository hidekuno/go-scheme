(() => {
    var history = document.querySelector('ul#history');
    var sExpression = document.querySelector('#sExpression');
    var url = new URL(location.href);
    var socket = new WebSocket(`ws://${url.host}/socket`);
    var escape = function(raw) {
        raw.replace(/&/g, "&amp;").replace(/</g, "&lt;").replace(/>/g, "&gt;").replace(/"/g, "&quot;").replace(/'/g, "&#039;");
        return raw;
    };
    var append = function(code) { history.innerHTML = `<li>${escape(code.text)}</li>` + history.innerHTML};
    var send = function () {
        message = sExpression.value;
        fetch("/message", {method:"POST",body:JSON.stringify({type:"EVAL",text: message}), headers:{"content-type":"application/json"}});
    };
    socket.onmessage = e => {
        code = JSON.parse(e.data);
        if (code.type == "CONNECT") {
            append(code);
        } else if (code.type == "LISPCODE") {
            sExpression.value = code.text;
        } else if (code.type == "EVAL") {
            append(code);
        }
    };
    socket.onerror = e => console.log("[ONERROR]", e);
    socket.onclose = e => console.log("[ONCLOSE]", e);

    document.querySelector('#evalButton').onclick = function() {
        eval();
        send();
    };
})();
