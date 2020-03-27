/*
   Go lang 3rd study program.
   This is prototype program mini scheme subset.

   hidekuno@gmail.com
*/
(() => {
    let history = document.querySelector('ul#history');
    let sExpression = document.querySelector('#sExpression');
    const url = new URL(location.href);
    const socket = new WebSocket(`ws://${url.host}/socket`);

    const escape = (raw) => {
        raw.replace(/&/g, "&amp;")
            .replace(/</g, "&lt;")
            .replace(/>/g, "&gt;")
            .replace(/"/g, "&quot;")
            .replace(/'/g, "&#039;");
        return raw;
    };
    const append = (code) => {history.innerHTML = `<li>${escape(code.text)}</li>` + history.innerHTML};
    const send = () => {
        message = sExpression.value;
        fetch("/message", {method: "POST",
                           body: JSON.stringify({type:"EVAL",text: message}),
                           headers: {"content-type":"application/json"}});
    };
    socket.onmessage = (e) => {
        code = JSON.parse(e.data);
        if (code.type == "CONNECT") {
            append(code);
        } else if (code.type == "LISPCODE") {
            sExpression.value = code.text;
        } else if (code.type == "EVAL") {
            append(code);
        }
    };
    socket.onerror = (e) => console.log("[ONERROR]", e);
    socket.onclose = (e) => console.log("[ONCLOSE]", e);

    document.querySelector('#evalButton').onclick = () => {
        eval();
        send();
    };
})();
