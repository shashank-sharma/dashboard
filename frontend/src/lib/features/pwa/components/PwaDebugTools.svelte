<script lang="ts">
  import { Button } from "$lib/components/ui/button";
  import {
    devControls,
    triggerInstallBanner,
    resetInstallBanner,
    toggleInstalledState,
    toggleOfflineMode,
    installPrompt,
    isPwaInstalled,
    isDebugModeEnabled,
  } from "$lib/features/pwa/services";
  import {
    Settings,
    X,
    Download,
    Smartphone,
    WifiOff,
    RefreshCcw,
  } from "lucide-svelte";
  import { fly } from "svelte/transition";
  import { onMount } from "svelte";

  let expanded = false;
  let promptAvailable = false;
  let isDebugMode = false;

  // Subscribe to the installPrompt store to know if a real prompt is available
  $: if ($installPrompt) {
    promptAvailable = true;
  }

  function togglePanel() {
    expanded = !expanded;
  }

  onMount(() => {
    isDebugMode = isDebugModeEnabled();
    console.log("Debug mode for debug tools: ", isDebugMode);
  });
</script>

{#if isDebugMode}
  <div class="fixed top-4 right-4 z-50">
    <Button
      variant="outline"
      size="icon"
      class="rounded-full shadow-lg bg-background"
      on:click={togglePanel}
      title="PWA Debug Tools"
    >
      {#if expanded}
        <X class="h-4 w-4" />
      {:else}
        <Settings class="h-4 w-4 text-amber-500" />
      {/if}
    </Button>

    <!-- Debug panel -->
    {#if expanded}
      <div
        class="absolute top-12 right-0 p-4 bg-card border rounded-lg shadow-lg w-80"
        transition:fly={{ y: 20, duration: 200 }}
      >
        <div class="flex items-center justify-between mb-3">
          <h3 class="font-medium">PWA Debug Tools</h3>
          <Button
            variant="ghost"
            size="sm"
            class="h-8 w-8 p-0"
            on:click={resetInstallBanner}
            title="Reset banner state (clear dismissed flag)"
          >
            <RefreshCcw class="h-3 w-3" />
          </Button>
        </div>

        <div class="space-y-4">
          <!-- PWA Install Status -->
          <div class="flex items-center justify-between">
            <span class="text-sm">
              Install Status:
              <span
                class={$isPwaInstalled || $devControls.simulateInstalled
                  ? "text-green-500"
                  : "text-amber-500"}
              >
                {$isPwaInstalled || $devControls.simulateInstalled
                  ? "Installed"
                  : "Not Installed"}
              </span>
            </span>

            <Button
              variant="outline"
              size="sm"
              class="gap-1"
              on:click={toggleInstalledState}
            >
              <Smartphone class="h-3 w-3" />
              Toggle Installed
            </Button>
          </div>

          <!-- PWA Prompt Status -->
          <div class="flex items-center justify-between">
            <span class="text-sm">
              Install Prompt:
              <span
                class={promptAvailable ? "text-green-500" : "text-amber-500"}
              >
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
              Show Banner
            </Button>
          </div>

          <!-- Offline Simulation -->
          <div class="flex items-center justify-between">
            <span class="text-sm">
              Network:
              <span
                class={$devControls.simulateOffline
                  ? "text-amber-500"
                  : "text-green-500"}
              >
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
              Toggle Offline
            </Button>
          </div>

          <div class="text-xs text-muted-foreground mt-2 pt-2 border-t">
            <p>If no install prompt appears:</p>
            <ul class="list-disc pl-4 mt-1 space-y-1">
              <li>Check browser console for errors</li>
              <li>Verify manifest.json is correctly configured</li>
              <li>
                Make sure your app is served over HTTPS (except on localhost)
              </li>
              <li>
                Chrome DevTools > Application has more PWA debugging tools
              </li>
            </ul>
          </div>
        </div>
      </div>
    {/if}
  </div>
{/if}
