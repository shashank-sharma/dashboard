<script lang="ts">
    import { onMount, setContext } from "svelte";
    import { navigating } from "$app/stores";
    import { Toaster } from "$lib/components/ui/sonner";
    import { fade } from "svelte/transition";
    import Loading from "$lib/components/Loading.svelte";
    import { theme } from "$lib/stores/theme.store";
    import "../app.css";
    import ThemeInitializer from "$lib/components/ThemeInitializer.svelte";
    import { initPwa } from "$lib/pwa";
    import PwaInstallBanner from "$lib/components/PwaInstallBanner.svelte";
    import PwaDebugTools from "$lib/components/PwaDebugTools.svelte";
    import { browser } from "$app/environment";
    let isLoading = true;

    setContext("theme", {
        theme,
        toggleTheme: theme.toggleTheme,
    });

    onMount(async () => {
        setTimeout(() => {
            isLoading = false;
        }, 500);

        // Initialize PWA functionality when on the client
        if (browser) {
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
        <PwaInstallBanner />
        
        <!-- PWA Debug Tools (only in development) -->
        {#if import.meta.env.DEV}
            <PwaDebugTools />
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
