/*
   Go lang 3rd study program.
   This is prototype program mini scheme subset.

   hidekuno@gmail.com
*/
(() => {
    if (!WebAssembly.instantiateStreaming) {
        WebAssembly.instantiateStreaming = async (resp, importObject) => {
            const source = await (await resp).arrayBuffer();
            return await WebAssembly.instantiate(source, importObject);
        };
    }

    const go = new Go();
    let mod, inst;
    WebAssembly.instantiateStreaming(fetch("wasm/lisp.wasm"), go.importObject).then(async (result) => {
        mod = result.module;
        inst = result.instance;
        console.clear();
        await go.run(inst);
    }).catch((err) => {
        console.error(err);
    });
})();
