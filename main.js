// Register service worker to control making site work offline
window.addEventListener('load', () => {
    if ('serviceWorker' in navigator) {
        navigator.serviceWorker
            .register('service-worker.js')
            .then(() => { console.log('[Service Worker] Registered'); })
            .catch((err) => { console.warn('[Service Worker] Error: ' + err); })
    }
});

// polyfill
if (!WebAssembly.instantiateStreaming) {
    WebAssembly.instantiateStreaming = async (resp, importObject) => {
        const source = await (await resp).arrayBuffer();
        return await WebAssembly.instantiate(source, importObject);
    };
}
const go = new Go();
WebAssembly
    .instantiateStreaming(fetch("app.wasm"), go.importObject)
    .then((result) => {
        go.run(result.instance);
    });
