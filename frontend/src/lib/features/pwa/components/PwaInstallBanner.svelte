<script lang="ts">
  import { onMount } from "svelte";
  import { fly } from "svelte/transition";
  import { X, Download, Info, Chrome, Menu, ArrowRight } from "lucide-svelte";
  import { Button } from "$lib/components/ui/button";
  import {
    installPrompt,
    isPwaInstalled,
    showInstallPrompt,
    promptUnavailableReason,
    isDebugModeEnabled,
    devControls,
  } from "$lib/features/pwa/services";
  import { browser } from "$app/environment";

  // Control when to show the banner
  let show = false;
  let installable = false;
  let dismissed = false;
  let showManualInstructions = false;
  let unavailabilityReason = "Unknown reason";
  let isDebugMode = false;
  let debug = {
    isInstalled: false,
    hasPrompt: false,
    isDismissed: false,
    isDevMode: false,
    reason: "",
  };

  // Browser detection for installation instructions
  let browserInfo = {
    isChrome: false,
    isFirefox: false,
    isSafari: false,
    isEdge: false,
    isIOS: false,
    isAndroid: false,
    name: "your browser",
  };

  const checkDismissed = () => {
    const sessionDismissed = sessionStorage.getItem("pwa-banner-dismissed");
    if (sessionDismissed === "true") {
      dismissed = true;
      debug.isDismissed = true;
    }
  };

  const dismiss = () => {
    show = false;
    dismissed = true;
    showManualInstructions = false;
    debug.isDismissed = true;
    sessionStorage.setItem("pwa-banner-dismissed", "true");
  };

  const forceShow = () => {
    dismissed = false;
    debug.isDismissed = false;
    sessionStorage.removeItem("pwa-banner-dismissed");
    updateShowState(true);
  };

  const handleInstall = async () => {
    const outcome = await showInstallPrompt();

    if (outcome === "unavailable") {
      // If no prompt is available, show manual instructions
      showManualInstructions = true;
    } else {
      dismiss();
    }
  };

  function detectBrowser() {
    if (!browser) return;

    const ua = navigator.userAgent;

    browserInfo.isIOS =
      /iPad|iPhone|iPod/.test(ua) ||
      (navigator.platform === "MacIntel" && navigator.maxTouchPoints > 1);
    browserInfo.isAndroid = /Android/.test(ua);

    if (/CriOS|Chrome/.test(ua) && !/Edge/.test(ua)) {
      browserInfo.isChrome = true;
      browserInfo.name = "Chrome";
    } else if (/Firefox/.test(ua)) {
      browserInfo.isFirefox = true;
      browserInfo.name = "Firefox";
    } else if (/Safari/.test(ua) && !/Chrome/.test(ua)) {
      browserInfo.isSafari = true;
      browserInfo.name = "Safari";
    } else if (/Edge/.test(ua)) {
      browserInfo.isEdge = true;
      browserInfo.name = "Edge";
    }
  }

  function updateShowState(force = false) {
    const isInstalled = $isPwaInstalled || $devControls.simulateInstalled;

    if (isInstalled) {
      show = false;
      return;
    }

    if (force) {
      show = true;
    } else {
      show = (installable || $devControls.showInstallPrompt) && !dismissed;
    }
  }

  function handleStorageEvent(e: StorageEvent) {
    if (e.key === "pwa-banner-dismissed") {
      checkDismissed();
      updateShowState();
    }
  }

  onMount(() => {
    if (!browser) return;

    // Check if debug mode is enabled
    isDebugMode = isDebugModeEnabled();
    debug.isDevMode = isDebugMode;

    checkDismissed();
    detectBrowser();

    const unsubscribePrompt = installPrompt.subscribe((value) => {
      installable = !!value;
      debug.hasPrompt = installable;
      updateShowState();
    });

    const unsubscribeInstalled = isPwaInstalled.subscribe((value) => {
      debug.isInstalled = value;
      updateShowState();
    });

    const unsubscribeReason = promptUnavailableReason.subscribe((value) => {
      if (value) {
        unavailabilityReason = value;
        debug.reason = value;
      }
    });

    const unsubscribeDev = devControls.subscribe((devState) => {
      updateShowState();
    });

    window.addEventListener("storage", handleStorageEvent);

    if (browser && window.location.search.includes("showPwaInstall=true")) {
      forceShow();
    }

    return () => {
      unsubscribePrompt();
      unsubscribeInstalled();
      unsubscribeDev();
      unsubscribeReason();
      window.removeEventListener("storage", handleStorageEvent);
    };
  });

  export { forceShow };
</script>

{#if show}
  <div
    class="fixed top-0 left-0 right-0 p-4 z-[100]"
    transition:fly={{ y: -100, duration: 300 }}
  >
    {#if showManualInstructions}
      <!-- Manual installation instructions -->
      <div
        class="bg-card text-card-foreground shadow-lg rounded-lg p-4 mx-auto max-w-lg border"
      >
        <div class="flex justify-between items-start mb-3">
          <div class="flex items-center gap-2">
            <Info class="h-5 w-5 text-blue-500" />
            <h3 class="font-semibold">Install Nen Space Manually</h3>
          </div>
          <Button variant="ghost" size="icon" on:click={dismiss}>
            <X class="w-4 h-4" />
          </Button>
        </div>

        <div class="text-sm space-y-3">
          <!-- Show the specific reason why automatic installation isn't available -->
          <div class="bg-muted/50 p-2 rounded text-muted-foreground mb-2">
            <p>
              Automatic installation is unavailable: <strong
                >{unavailabilityReason}</strong
              >
            </p>
          </div>

          {#if browserInfo.isIOS}
            <!-- iOS Safari instructions -->
            <p>To install Nen Space on your iPhone or iPad:</p>
            <ol class="list-decimal space-y-2 pl-5">
              <li>
                Tap the <span class="inline-block px-2 font-semibold"
                  >Share</span
                > button in Safari
              </li>
              <li>
                Scroll down and tap <span
                  class="inline-block px-2 font-semibold"
                  >Add to Home Screen</span
                >
              </li>
              <li>
                Tap <span class="inline-block px-2 font-semibold">Add</span>
              </li>
            </ol>
          {:else if browserInfo.isChrome}
            <!-- Chrome instructions -->
            <p>To install Nen Space in Chrome:</p>
            <ol class="list-decimal space-y-2 pl-5">
              <li>
                Click the menu button <span
                  class="inline-block px-2 font-semibold"
                  ><Menu class="h-4 w-4 inline" /></span
                > in the top right
              </li>
              <li>
                Select <span class="inline-block px-2 font-semibold"
                  >Install App</span
                >
                or
                <span class="inline-block px-2 font-semibold"
                  >Install Nen Space</span
                >
              </li>
            </ol>
          {:else if browserInfo.isFirefox}
            <!-- Firefox instructions -->
            <p>To install Nen Space in Firefox:</p>
            <ol class="list-decimal space-y-2 pl-5">
              <li>
                Click the menu button <span
                  class="inline-block px-2 font-semibold"
                  ><Menu class="h-4 w-4 inline" /></span
                > in the top right
              </li>
              <li>
                Select <span class="inline-block px-2 font-semibold"
                  >Install App</span
                >
                or
                <span class="inline-block px-2 font-semibold"
                  >Add to Home Screen</span
                >
              </li>
            </ol>
          {:else}
            <!-- Generic instructions -->
            <p>To install Nen Space on your device:</p>
            <ol class="list-decimal space-y-2 pl-5">
              <li>
                Open the menu in {browserInfo.name}
              </li>
              <li>
                Look for an option like <span
                  class="inline-block px-2 font-semibold">Install App</span
                >,
                <span class="inline-block px-2 font-semibold"
                  >Add to Home Screen</span
                >, or
                <span class="inline-block px-2 font-semibold"
                  >Install Nen Space</span
                >
              </li>
            </ol>
          {/if}

          <div class="pt-2 border-t mt-2">
            <p class="text-muted-foreground text-xs">
              After installation, you can launch the app from your home
              screen/desktop and use it offline.
            </p>
          </div>
        </div>
      </div>
    {:else}
      <!-- Main banner content -->
      <div
        class="bg-card text-card-foreground shadow-lg rounded-lg p-4 mx-auto max-w-lg border"
      >
        <div class="flex justify-between items-start mb-3">
          <div class="flex items-center gap-2">
            <Download class="h-5 w-5 text-primary" />
            <h3 class="font-semibold">Install Nen Space</h3>
          </div>
          <Button variant="ghost" size="icon" on:click={dismiss}>
            <X class="w-4 h-4" />
          </Button>
        </div>

        <p class="text-sm mb-3">
          Install Nen Space on your device for a better experience with offline
          access and faster loading.
        </p>

        <div class="flex justify-end gap-2">
          <Button variant="outline" size="sm" on:click={dismiss}>
            Not now
          </Button>
          <Button variant="default" size="sm" on:click={handleInstall}>
            <Download class="w-4 h-4 mr-2" /> Install
          </Button>
        </div>
      </div>
    {/if}
  </div>
{/if}

<!-- Debug UI - only visible in debug mode -->
{#if isDebugMode}
  <div class="fixed bottom-20 right-4 z-50">
    <Button variant="outline" size="sm" on:click={forceShow}>
      Test PWA Banner
    </Button>
    <div class="mt-2 text-xs bg-card p-2 rounded border shadow-sm">
      <div><strong>Debug:</strong></div>
      <div>Installed: {debug.isInstalled ? "✅" : "❌"}</div>
      <div>Has Prompt: {debug.hasPrompt ? "✅" : "❌"}</div>
      <div>Dismissed: {debug.isDismissed ? "✅" : "❌"}</div>
      <div>Debug Mode: {debug.isDevMode ? "✅" : "❌"}</div>
      <div>Browser: {browserInfo.name}</div>
      {#if debug.reason}
        <div class="text-amber-500">Reason: {debug.reason}</div>
      {/if}
    </div>
  </div>
{/if}
