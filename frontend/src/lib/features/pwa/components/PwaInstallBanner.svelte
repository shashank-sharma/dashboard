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
      console.log("[PWA] Banner previously dismissed in this session");
    } else {
      console.log("[PWA] Banner not previously dismissed");
    }
  };

  const dismiss = () => {
    show = false;
    dismissed = true;
    showManualInstructions = false;
    debug.isDismissed = true;
    sessionStorage.setItem("pwa-banner-dismissed", "true");
    console.log("[PWA] Banner dismissed and saved to session storage");
  };

  const forceShow = () => {
    dismissed = false;
    debug.isDismissed = false;
    sessionStorage.removeItem("pwa-banner-dismissed");
    console.log("[PWA] Force showing banner, dismissal status reset");
    updateShowState(true);
  };

  const handleInstall = async () => {
    console.log("[PWA] Install button clicked");
    const outcome = await showInstallPrompt();
    console.log("[PWA] Install prompt outcome:", outcome);

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
      console.log("[PWA] App is installed, hiding banner");
      return;
    }

    if (force) {
      show = true;
      console.log("[PWA] Showing banner (forced)");
    } else {
      show = (installable || $devControls.showInstallPrompt) && !dismissed;
      console.log("[PWA] Banner visibility:", show, {
        installable,
        devShowPrompt: $devControls.showInstallPrompt,
        isInstalled,
        dismissed,
      });
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

    console.log("[PWA] Checking installation criteria...");

    if (isDebugMode) {
      console.log("[PWA-DEV] Debug mode detected");
    }

    const unsubscribePrompt = installPrompt.subscribe((value) => {
      installable = !!value;
      debug.hasPrompt = installable;
      console.log("[PWA] Install prompt available:", installable);
      updateShowState();
    });

    const unsubscribeInstalled = isPwaInstalled.subscribe((value) => {
      debug.isInstalled = value;
      console.log("[PWA] Is installed:", value);
      updateShowState();
    });

    const unsubscribeReason = promptUnavailableReason.subscribe((value) => {
      if (value) {
        unavailabilityReason = value;
        debug.reason = value;
        console.log("[PWA] Prompt unavailable reason:", value);
      }
    });

    const unsubscribeDev = devControls.subscribe((devState) => {
      if (isDebugMode) {
        console.log("[PWA-DEV] Dev controls state:", devState);
      }
      updateShowState();
    });

    window.addEventListener("storage", handleStorageEvent);

    if (isDebugMode && !installable) {
      console.log(
        "[PWA-DEV] No install prompt available. Try using PwaDebugTools to trigger manually.",
      );
    }

    if (browser && window.location.search.includes("showPwaInstall=true")) {
      console.log("[PWA] Detected showPwaInstall parameter in URL");
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
            <h3 class="font-semibold">Install Manually</h3>
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
            <p>To install this app on your iPhone or iPad:</p>
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
                Tap <span class="inline-block px-2 font-semibold">Add</span> in the
                top right
              </li>
            </ol>
          {:else if browserInfo.isAndroid && browserInfo.isChrome}
            <!-- Android Chrome instructions -->
            <p>To install this app on your Android device:</p>
            <ol class="list-decimal space-y-2 pl-5">
              <li>
                Tap the <span class="inline-block px-2 font-semibold"
                  >Three dots</span
                > menu in Chrome
              </li>
              <li>
                Tap <span class="inline-block px-2 font-semibold"
                  >Install app</span
                >
                or
                <span class="inline-block px-2 font-semibold"
                  >Add to Home screen</span
                >
              </li>
              <li>Follow the on-screen instructions</li>
            </ol>
          {:else if browserInfo.isChrome || browserInfo.isEdge}
            <!-- Chrome/Edge desktop instructions -->
            <p>To install this app in {browserInfo.name}:</p>
            <ol class="list-decimal space-y-2 pl-5">
              <li>
                Click the <span class="inline-block px-2 font-semibold"
                  >Install</span
                >
                icon in the address bar <Chrome class="h-4 w-4 inline" />
              </li>
              <li>
                Or, click the <span class="inline-block px-2 font-semibold"
                  >Three dots</span
                > menu
              </li>
              <li>
                Select <span class="inline-block px-2 font-semibold"
                  >Install [App Name]...</span
                >
              </li>
            </ol>
          {:else if browserInfo.isFirefox}
            <!-- Firefox instructions -->
            <p>To install this app in Firefox:</p>
            <ol class="list-decimal space-y-2 pl-5">
              <li>
                Click the <span class="inline-block px-2 font-semibold"
                  >Three dots</span
                >
                menu in the address bar <Menu class="h-4 w-4 inline" />
              </li>
              <li>
                Select <span class="inline-block px-2 font-semibold"
                  >Install</span
                >
              </li>
            </ol>
          {:else}
            <!-- Generic instructions -->
            <p>To install this app:</p>
            <ol class="list-decimal space-y-2 pl-5">
              <li>
                Look for an <span class="inline-block px-2 font-semibold"
                  >Install</span
                >
                or
                <span class="inline-block px-2 font-semibold"
                  >Add to Home Screen</span
                > option in your browser's menu
              </li>
              <li>
                On mobile, this is usually in the share menu or browser options
              </li>
              <li>
                On desktop, look for an icon in the address bar or in the main
                browser menu
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
      <!-- Standard install banner -->
      <div
        class="bg-card text-card-foreground shadow-lg rounded-lg p-4 mx-auto max-w-md flex items-center justify-between border"
      >
        <div class="flex items-center space-x-4">
          <div
            class="w-10 h-10 bg-primary/10 rounded-full flex items-center justify-center"
          >
            <Download class="w-5 h-5 text-primary" />
          </div>
          <div>
            <h3 class="font-semibold">Install App</h3>
            <p class="text-sm text-muted-foreground">
              Add to your home screen for a better experience
            </p>
          </div>
        </div>

        <div class="flex space-x-2">
          <Button variant="ghost" size="icon" on:click={dismiss}>
            <X class="w-4 h-4" />
          </Button>
          <Button variant="default" on:click={handleInstall}>Install</Button>
        </div>
      </div>
    {/if}
  </div>
{/if}

<!-- Debug UI - only visible in debug mode -->
{#if isDebugMode}
  <div class="fixed bottom-4 right-4 z-50">
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
