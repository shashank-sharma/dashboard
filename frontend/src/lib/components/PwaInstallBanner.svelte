<script lang="ts">
    import { onMount } from "svelte";
    import { fly } from "svelte/transition";
    import { X, Download } from "lucide-svelte";
    import { Button } from "$lib/components/ui/button";
    import { installPrompt, isPwaInstalled, showInstallPrompt } from "$lib/pwa";
    import { browser } from '$app/environment';

    // Import dev controls in development mode only
    import { devControls } from "$lib/pwa-dev";

    // Control when to show the banner
    let show = false;
    let installable = false;
    let dismissed = false;

    // Check if we've already dismissed the banner in this session
    const checkDismissed = () => {
        const sessionDismissed = sessionStorage.getItem("pwa-banner-dismissed");
        if (sessionDismissed === "true") {
            dismissed = true;
        }
    };

    // Dismiss the banner and save to session storage
    const dismiss = () => {
        show = false;
        dismissed = true;
        sessionStorage.setItem("pwa-banner-dismissed", "true");
    };

    // Handle install click
    const handleInstall = async () => {
        await showInstallPrompt();
        dismiss();
    };

    onMount(() => {
        if (!browser) return;
    
    checkDismissed();
    
    // Subscribe to the installPrompt store
    const unsubscribePrompt = installPrompt.subscribe(value => {
      installable = !!value;
      updateShowState();
    });
    
    // In development, also subscribe to the dev controls
    const unsubscribeDev = devControls.subscribe(() => {
      updateShowState();
    });
    
    // Function to determine whether to show the banner
    function updateShowState() {
      const isInstalled = $isPwaInstalled || $devControls.simulateInstalled;
      
      if (isInstalled) {
        dismiss();
        return;
      }
      
      // Show if we have a real prompt or if we're forcing it in dev mode
      show = (installable || $devControls.showInstallPrompt) && !isInstalled && !dismissed;
    }
    
    return () => {
      unsubscribePrompt();
      unsubscribeDev();
    };
    });
</script>

{#if show}
  <div 
    class="fixed bottom-0 left-0 right-0 p-4 z-50"
    transition:fly={{ y: 100, duration: 300 }}
  >
    <div class="bg-card text-card-foreground shadow-lg rounded-lg p-4 mx-auto max-w-md flex items-center justify-between border">
      <div class="flex items-center space-x-4">
        <div class="w-10 h-10 bg-primary/10 rounded-full flex items-center justify-center">
          <Download class="w-5 h-5 text-primary" />
        </div>
        <div>
          <h3 class="font-semibold">Install App</h3>
          <p class="text-sm text-muted-foreground">Add to your home screen for a better experience</p>
        </div>
      </div>
      
      <div class="flex space-x-2">
        <Button variant="ghost" size="icon" on:click={dismiss}>
          <X class="w-4 h-4" />
        </Button>
        <Button variant="default" on:click={handleInstall}>
          Install
        </Button>
      </div>
    </div>
  </div>
{/if}