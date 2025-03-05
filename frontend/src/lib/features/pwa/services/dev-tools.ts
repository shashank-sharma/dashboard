// PWA development tools for testing and debugging
import { writable, get } from 'svelte/store';
import { browser } from '$app/environment';
import { installPrompt } from './core';

export const devControls = writable({
  showInstallPrompt: false,
  simulateInstalled: false,
  simulateOffline: false
});

export function triggerInstallBanner() {
  if (!browser) return;
  
  console.log('[PWA-DEV] Manually triggering install banner');
  
  const event = new CustomEvent(
    "force-show-pwa-banner",
);
window.dispatchEvent(event);
}

export function resetInstallBanner() {
  if (!browser) return;
  
  console.log('[PWA-DEV] Resetting install banner state');
  sessionStorage.removeItem("pwa-banner-dismissed");
  devControls.update(state => ({ ...state, showInstallPrompt: false }));
}

export function toggleInstalledState() {
  if (!browser) return;
  
  const currentState = get(devControls);
  const newValue = !currentState.simulateInstalled;
  
  devControls.update(state => ({ 
    ...state, 
    simulateInstalled: newValue
  }));
  
  console.log('[PWA-DEV] Simulating installed state:', newValue);
}

export function toggleOfflineMode() {
  if (!browser) return;
  
  const currentState = get(devControls);
  const newValue = !currentState.simulateOffline;
  
  devControls.update(state => ({ 
    ...state, 
    simulateOffline: newValue 
  }));
  
  console.log('[PWA-DEV] Simulating offline mode:', newValue);
  
  if ('serviceWorker' in navigator) {
    const controller = navigator.serviceWorker.controller;
    if (controller) {
      controller.postMessage({
        type: 'OFFLINE_MODE',
        offline: newValue
      });
    }
  }
} 