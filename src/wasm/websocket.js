(() => {
    var sExpression = document.querySelector('#sExpression');
    var url = new URL(location.href);
    var socket = new WebSocket(`ws://${url.host}/socket`);

    socket.onmessage = e => { code = JSON.parse(e.data); sExpression.value = code.text; };
    socket.onerror = e => console.log("[ONERROR]", e);
    socket.onclose = e => console.log("[ONCLOSE]", e);
})();
