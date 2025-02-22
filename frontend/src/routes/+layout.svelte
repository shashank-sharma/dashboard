<script lang="ts">
    import { onMount } from "svelte";
    import { navigating } from "$app/stores";
    import { Toaster } from "$lib/components/ui/sonner";
    import { fade } from "svelte/transition";
    import Loading from "$lib/components/Loading.svelte";
    import "../app.css";
    let isLoading = true;

    onMount(async () => {
        setTimeout(() => {
            isLoading = false;
        }, 500);
    });
</script>

{#if isLoading || $navigating}
    <Loading />
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
