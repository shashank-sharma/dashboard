import { readable, derived, type Readable } from 'svelte/store';
import { browser } from '$app/environment';

const MOBILE_BREAKPOINT = 768;


function createResponsiveStore() {
  const windowSize = readable({ width: browser ? window.innerWidth : 1024, height: browser ? window.innerHeight : 768 }, (set) => {
    if (!browser) return;
    
    // Handler for resize events
    const handleResize = () => {
      set({
        width: window.innerWidth,
        height: window.innerHeight
      });
    };
    
    handleResize();
    window.addEventListener('resize', handleResize);
    return () => {
      window.removeEventListener('resize', handleResize);
    };
  });
  
  const isMobile = derived(windowSize, ($size) => $size.width < MOBILE_BREAKPOINT);
  const isTablet = derived(windowSize, ($size) => $size.width >= MOBILE_BREAKPOINT && $size.width < 1024);
  const isDesktop = derived(windowSize, ($size) => $size.width >= 1024);
  const orientation = derived(windowSize, ($size) => 
    $size.width >= $size.height ? 'landscape' : 'portrait'
  );
  
  return {
    windowSize,
    isMobile,
    isTablet, 
    isDesktop,
    orientation
  };
}

export const device = createResponsiveStore();

export const isMobile = device.isMobile;
export const isTablet = device.isTablet;
export const isDesktop = device.isDesktop;
export const orientation = device.orientation;