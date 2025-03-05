<script lang="ts">
    import { Button } from '$lib/components/ui/button';
    import { devControls, triggerInstallBanner, toggleInstalledState, toggleOfflineMode } from '$lib/pwa-dev';
    import { installPrompt, isPwaInstalled } from '$lib/pwa';
    import { Settings, X, Download, Smartphone, WifiOff } from 'lucide-svelte';
    import { fly } from 'svelte/transition';
    
    let expanded = false;
    let promptAvailable = false;
    
    // Subscribe to the installPrompt store to know if a real prompt is available
    $: if ($installPrompt) {
      promptAvailable = true;
    }
    
    function togglePanel() {
      expanded = !expanded;
    }
  </script>
  
  {#if import.meta.env.DEV}
    <div class="fixed bottom-4 right-4 z-50">
      <!-- Toggle button -->
      <Button 
        variant="outline" 
        size="icon" 
        class="rounded-full shadow-lg bg-background"
        on:click={togglePanel}
      >
        {#if expanded}
          <X class="h-4 w-4" />
        {:else}
          <Settings class="h-4 w-4" />
        {/if}
      </Button>
      
      <!-- Debug panel -->
      {#if expanded}
        <div 
          class="absolute bottom-12 right-0 p-4 bg-card border rounded-lg shadow-lg w-72"
          transition:fly={{ y: 20, duration: 200 }}
        >
          <h3 class="font-medium mb-3">PWA Debug Tools</h3>
          
          <div class="space-y-3">
            <!-- PWA Install Status -->
            <div class="flex items-center justify-between">
              <span class="text-sm">
                Install Status:
                <span class={$isPwaInstalled || $devControls.simulateInstalled ? "text-green-500" : "text-amber-500"}>
                  {$isPwaInstalled || $devControls.simulateInstalled ? "Installed" : "Not Installed"}
                </span>
              </span>
              
              <Button 
                variant="outline" 
                size="sm" 
                class="gap-1"
                on:click={toggleInstalledState}
              >
                <Smartphone class="h-3 w-3" />
                Simulate
              </Button>
            </div>
            
            <!-- PWA Prompt Status -->
            <div class="flex items-center justify-between">
              <span class="text-sm">
                Install Prompt:
                <span class={promptAvailable ? "text-green-500" : "text-amber-500"}>
                  {promptAvailable ? "Available" : "Not Available"}
                </span>
              </span>
              
              <Button 
                variant="outline" 
                size="sm" 
                class="gap-1"
                on:click={triggerInstallBanner}
              >
                <Download class="h-3 w-3" />
                Test Banner
              </Button>
            </div>
            
            <!-- Offline Simulation -->
            <div class="flex items-center justify-between">
              <span class="text-sm">
                Network:
                <span class={$devControls.simulateOffline ? "text-amber-500" : "text-green-500"}>
                  {$devControls.simulateOffline ? "Offline" : "Online"}
                </span>
              </span>
              
              <Button 
                variant="outline" 
                size="sm" 
                class="gap-1"
                on:click={toggleOfflineMode}
              >
                <WifiOff class="h-3 w-3" />
                Toggle
              </Button>
            </div>
            
            <div class="text-xs text-muted-foreground mt-2 pt-2 border-t">
              Remember to check Chrome DevTools > Application for more PWA debugging options.
            </div>
          </div>
        </div>
      {/if}
    </div>
  {/if}