import { writable } from 'svelte/store';

// Store for the installation prompt event
export const installPrompt = writable<Event | null>(null);

// Store for the PWA installation status
export const isPwaInstalled = writable<boolean>(false);

// Check if the app is already installed
export const checkInstallStatus = () => {
  if (window.matchMedia('(display-mode: standalone)').matches) {
    isPwaInstalled.set(true);
  }
};

// Register service worker
export const registerServiceWorker = async () => {
  if ('serviceWorker' in navigator) {
    try {
      await navigator.serviceWorker.register('/sw.js');
      console.log('Service worker registered successfully');
    } catch (error) {
      console.error('Service worker registration failed:', error);
    }
  }
};

// Initialize PWA functionality
export const initPwa = () => {
  // Register service worker
  registerServiceWorker();
  
  // Check if already installed
  checkInstallStatus();
  
  // Listen for the beforeinstallprompt event
  window.addEventListener('beforeinstallprompt', (e) => {
    // Prevent Chrome 67 and earlier from automatically showing the prompt
    e.preventDefault();
    // Store the event so it can be triggered later
    installPrompt.set(e);
  });
  
  // Listen for app installed event
  window.addEventListener('appinstalled', () => {
    // Clear the prompt once installed
    installPrompt.set(null);
    isPwaInstalled.set(true);
    console.log('PWA was installed');
  });
};

// Prompt the user to install the PWA
export const showInstallPrompt = async () => {
  let promptEvent: any;
  
  // Get the stored event from the store
  installPrompt.subscribe(value => {
    promptEvent = value;
  })();
  
  if (promptEvent) {
    // Show the prompt
    promptEvent.prompt();
    
    // Wait for the user to respond to the prompt
    const userChoice = await promptEvent.userChoice;
    
    // Check the user's response
    if (userChoice.outcome === 'accepted') {
      console.log('User accepted the install prompt');
    } else {
      console.log('User dismissed the install prompt');
    }
    
    // Clear the prompt
    installPrompt.set(null);
  }
};