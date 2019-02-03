(() => {
    var url = new URL(location.href);
    var socket = new WebSocket(`ws//${url.host}/socket`);
    socket.onmessage = e => {sExpression.value = JSON.parse(e.data);};
    socket.onerror = e => console.log("[ONERROR]", e);
    socket.onclose = e => console.log("[ONCLOSE]", e);
})();
