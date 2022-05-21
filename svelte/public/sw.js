self.addEventListener('install', e => {
  e.waitUntil( async () => {
    const cache = caches.open('offpost')
    await cache.addAll(['/index.html'])
  })
})

self.addEventListener('fetch', (e) => {
  console.log(`tried to request ${e.request.status}`)
});