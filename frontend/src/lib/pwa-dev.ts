// src/lib/pwa-dev.ts
// This file provides development helpers for PWA testing

import { writable } from 'svelte/store';
import { browser } from '$app/environment';

// Store for manual control of PWA features during development
export const devControls = writable({
  showInstallPrompt: false,
  simulateInstalled: false,
  simulateOffline: false
});

// Manually trigger the install banner for testing
export function triggerInstallBanner() {
  if (browser) {
    devControls.update(state => ({ ...state, showInstallPrompt: true }));
  }
}

// Simulate the PWA being installed
export function toggleInstalledState() {
  if (browser) {
    devControls.update(state => ({ 
      ...state, 
      simulateInstalled: !state.simulateInstalled 
    }));
  }
}

// Simulate offline mode
export function toggleOfflineMode() {
  if (browser) {
    devControls.update(state => ({ 
      ...state, 
      simulateOffline: !state.simulateOffline 
    }));
    
    // Actually apply the offline simulation if supported
    if ('serviceWorker' in navigator) {
      const controller = navigator.serviceWorker.controller;
      if (controller) {
        controller.postMessage({
          type: 'OFFLINE_MODE',
          offline: !devControls.simulateOffline
        });
      }
    }
  }
}