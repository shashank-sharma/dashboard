// Service worker version - increment this when you make changes
const CACHE_VERSION = 'v3';
const CACHE_NAME = `app-cache-${CACHE_VERSION}`;
// const DEV_MODE = self.location.hostname === 'localhost' || self.location.hostname === '127.0.0.1';
const DEV_MODE = false;

// Debug flags for testing
let simulateOffline = false;
let logCacheOperations = DEV_MODE;

// Assets to cache
const ASSETS_TO_CACHE = [
  '/',
  '/manifest.json',
  '/favicon.png',
  '/icons/apple-icon-180.png',
  '/icons/manifest-icon-192.maskable.png',
  '/icons/manifest-icon-512.maskable.png',
  '/fonts/Gilroy.woff2'
];

// Helper for logging in development
const devLog = (...args) => {
  if (DEV_MODE) {
    console.log('[Service Worker]', ...args);
  }
};

// Listen for messages from the main thread
self.addEventListener('message', (event) => {
  devLog('Message received', event.data);
  
  if (event.data && event.data.type === 'OFFLINE_MODE') {
    simulateOffline = event.data.offline;
    devLog('Offline mode set to:', simulateOffline);
  }
  
  if (event.data && event.data.type === 'SKIP_WAITING') {
    self.skipWaiting();
  }
});

// Install event - cache critical assets
self.addEventListener('install', (event) => {
  devLog('Installing service worker...');
  
  event.waitUntil(
    caches.open(CACHE_NAME)
      .then((cache) => {
        devLog('Caching app shell...');
        return cache.addAll(ASSETS_TO_CACHE);
      })
      .then(() => {
        devLog('Service worker installed');
        return self.skipWaiting();
      })
      .catch(error => {
        devLog('Installation failed:', error);
      })
  );
});

// Activate event - clean up old caches
self.addEventListener('activate', (event) => {
  devLog('Activating service worker...');
  const cacheWhitelist = [CACHE_NAME];
  
  event.waitUntil(
    caches.keys().then((cacheNames) => {
      return Promise.all(
        cacheNames.map((cacheName) => {
          if (cacheWhitelist.indexOf(cacheName) === -1) {
            devLog('Deleting old cache:', cacheName);
            return caches.delete(cacheName);
          }
        })
      );
    }).then(() => {
      devLog('Service worker activated and controlling page');
      return self.clients.claim();
    })
  );
});

// Fetch event - serve from cache or network
self.addEventListener('fetch', (event) => {
  if (event.request.method !== 'GET') return;  
  const url = new URL(event.request.url);

  if (
    url.pathname.startsWith('/browser-sync/') ||
    url.pathname.startsWith('/__vite') ||
    url.pathname.startsWith('/node_modules/')
  ) {
    return;
  }
  
  // For development debugging
  if (DEV_MODE && logCacheOperations) {
    devLog('Fetch event for:', event.request.url);
  }
  
  event.respondWith(
    (async () => {
      const cachedResponse = await caches.match(event.request);
      if (cachedResponse) {
        if (DEV_MODE && logCacheOperations) {
          devLog('Found in cache:', event.request.url);
        }
        return cachedResponse;
      }
      
      // If simulating offline mode in development and not found in cache, return offline response
      if (DEV_MODE && simulateOffline) {
        devLog('Simulating offline mode, returning offline response');
        
        if (event.request.destination === 'document') {
          return caches.match('/offline.html').catch(() => {
            return new Response('<h1>You are offline</h1><p>This is a simulated offline response.</p>', {
              headers: { 'Content-Type': 'text/html' }
            });
          });
        }
        
        return new Response('Offline Mode', { 
          status: 503, 
          statusText: 'Simulated Offline Mode' 
        });
      }
      
      // Otherwise, try network
      try {
        const networkResponse = await fetch(event.request);
        
        // Check if we should cache this response
        if (
          networkResponse && 
          networkResponse.status === 200 &&
          networkResponse.type === 'basic' &&
          !event.request.url.includes('/api/') // Don't cache API requests
        ) {
          const cache = await caches.open(CACHE_NAME);
          if (DEV_MODE && logCacheOperations) {
            devLog('Caching new resource:', event.request.url);
          }
          cache.put(event.request, networkResponse.clone());
        }
        
        return networkResponse;
      } catch (error) {
        // Network request failed and resource not in cache
        if (DEV_MODE) {
          devLog('Network request failed:', error);
          
          if (event.request.destination === 'document') {
            return caches.match('/offline.html').catch(() => {
              return new Response('<h1>You are offline</h1><p>The requested page is not available offline.</p>', {
                headers: { 'Content-Type': 'text/html' }
              });
            });
          }
        }
        
        // Return error for other resources
        return new Response('Network error occurred', { 
          status: 503, 
          statusText: 'Service Unavailable' 
        });
      }
    })()
  );
});