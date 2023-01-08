/*
Sources:
- https://developers.google.com/web/fundamentals/primers/service-workers/
- https://github.com/mdn/pwa-examples
- https://developer.mozilla.org/es/docs/Web/Progressive_web_apps
*/

var cacheName = 'dice-roller-v1';

self.addEventListener('install', (event) => {
    console.log('[Service Worker] Install');
    event.waitUntil(
        caches.open(cacheName).then((cache) => cache.addAll([
            '/',
            '/app.wasm',
            '/index.html',
            '/main.js',
            '/manifest.json',
            '/style.css',
            '/wasm_exec.js',
        ])),
    );
});

self.addEventListener('fetch', (event) => {
    event.respondWith(
        caches.match(event.request).then((response) => {
            console.log('[Service Worker] Fetching: ' + event.request.url);
            return response || fetch(event.request).then((response) => {
                return caches.open(cacheName).then((cache) => {
                    console.log('[Service Worker] Catching: ' + event.request.url);
                    cache.put(event.request, response.clone());
                    return response;
                });
            });
        })
    );
});
