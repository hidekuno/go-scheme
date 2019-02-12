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

    socket.onmessage = e => { code = JSON.parse(e.data); sExpression.value = code.text; append(code);};
    socket.onerror = e => console.log("[ONERROR]", e);
    socket.onclose = e => console.log("[ONCLOSE]", e);
})();
