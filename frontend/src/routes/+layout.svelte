<script lang="ts">
    import { onMount, setContext } from "svelte";
    import { navigating } from "$app/stores";
    import { Toaster } from "$lib/components/ui/sonner";
    import { fade } from "svelte/transition";
    import Loading from "$lib/components/Loading.svelte";
    import { theme } from "$lib/stores/theme.store";
    import "../app.css";
    import ThemeInitializer from "$lib/components/ThemeInitializer.svelte";
    import {
        initPwa,
        isDebugModeEnabled,
        disableServiceWorkerCaching,
        enableServiceWorkerCaching,
        isServiceWorkerCachingDisabled,
    } from "$lib/features/pwa/services";
    import { Button } from "$lib/components/ui/button";
    import { RefreshCw } from "lucide-svelte";
    import { browser } from "$app/environment";
    import UpdateDetector from "$lib/features/pwa/components/UpdateDetector.svelte";
    import { toast } from "svelte-sonner";
    let isLoading = true;
    let isDebugMode = false;
    let isCachingDisabled = false;

    setContext("theme", {
        theme,
        toggleTheme: theme.toggleTheme,
    });

    async function handleFreshReload() {
        toast.loading("Clearing caches and reloading...");
        // Force clear caches - this is more thorough than just disabling caching
        try {
            if (browser && "caches" in window) {
                const cacheKeys = await window.caches.keys();
                await Promise.all(
                    cacheKeys.map((key) => window.caches.delete(key)),
                );

                await disableServiceWorkerCaching();
                isCachingDisabled = true;

                const timestamp = new Date().getTime();
                const url = new URL(window.location.href);
                url.searchParams.set("_t", timestamp.toString());

                window.location.href = url.toString();
            } else {
                window.location.reload();
            }
        } catch (error) {
            console.error("[PWA] Error during cache clearing:", error);
            window.location.reload();
        }
    }

    async function toggleCaching() {
        if (isCachingDisabled) {
            toast.loading("Enabling caching...");
            const success = await enableServiceWorkerCaching();
            if (success) {
                toast.success("Caching enabled");
                isCachingDisabled = false;
            } else {
                toast.error("Failed to enable caching");
            }
        } else {
            toast.loading("Disabling caching...");
            const success = await disableServiceWorkerCaching();
            if (success) {
                toast.success(
                    "Caching disabled - all content will be fetched fresh",
                );
                isCachingDisabled = true;
            } else {
                toast.error("Failed to disable caching");
            }
        }
    }

    onMount(async () => {
        setTimeout(() => {
            isLoading = false;
        }, 500);

        if (browser) {
            isDebugMode = isDebugModeEnabled();

            if (isDebugMode) {
                try {
                    isCachingDisabled = await isServiceWorkerCachingDisabled();
                    console.log(
                        "[PWA] Service Worker caching is disabled:",
                        isCachingDisabled,
                    );
                } catch (error) {
                    console.error(
                        "[PWA] Error checking caching status:",
                        error,
                    );
                }
            }

            initPwa();
        }
    });
</script>

<ThemeInitializer />

{#if isLoading || $navigating}
    <Loading />
{/if}

<!-- PWA Install Banner -->
{#if browser}
    <!-- Service Worker Update Detector -->
    <UpdateDetector />

    <!-- Add cache toggle button to the dev tools section -->
    {#if isDebugMode}
        <div class="fixed top-4 left-4 z-50 flex gap-2">
            <Button
                variant="destructive"
                size="sm"
                class="shadow-lg gap-1"
                on:click={handleFreshReload}
            >
                <RefreshCw class="h-4 w-4" /> Force Refresh
            </Button>

            <!-- Add a new button to toggle caching -->
            <Button
                variant={isCachingDisabled ? "default" : "outline"}
                size="sm"
                class="shadow-lg gap-1"
                on:click={toggleCaching}
            >
                {#if isCachingDisabled}
                    <span>Caching Disabled</span>
                {:else}
                    <span>Disable Caching</span>
                {/if}
            </Button>
        </div>
    {/if}
{/if}

<div
    class="app-container"
    in:fade={{ duration: 300 }}
    class:pointer-events-none={$navigating}
    class:opacity-50={$navigating}
>
    <Toaster />

    <main class="">
        <slot />
    </main>
</div>

<style>
    @font-face {
        font-family: "Gilroy";
        src: url("/fonts/Gilroy.woff2");
    }

    :global(html) {
        height: 100%;
        margin: 0;
        padding: 0;
        font-family: Gilroy, serif;
    }
    .app-container {
        min-height: 100vh;
        background-color: hsl(var(--background));
        color: hsl(var(--foreground));
    }
</style>
