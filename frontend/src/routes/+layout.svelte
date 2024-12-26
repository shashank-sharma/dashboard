<script lang="ts">
    import { onMount } from "svelte";
    import SettingsModal from "$lib/components/SettingsModal.svelte";
    import { navigating } from "$app/stores";
    import { Toaster } from "$lib/components/ui/sonner";
    import { fade } from "svelte/transition";
    import Loading from "$lib/components/Loading.svelte";
    import { pb, refreshToken } from "$lib/pocketbase";
    import { goto } from "$app/navigation";
    import "../app.css";
    import { Button } from "$lib/components/ui/button";
    import { Sun, Moon } from "lucide-svelte";
    import { writable } from "svelte/store";
    import { setContext } from "svelte";

    let theme = writable("light");
    let isLoading = true;

    function toggleTheme() {
        theme.update((t) => (t === "dark" ? "light" : "dark"));
    }

    setContext("theme", {
        theme: theme,
        toggleTheme: toggleTheme,
    });

    onMount(async () => {
        try {
            if (pb.authStore.isValid) {
                await refreshToken();
            } else if (!window.location.pathname.startsWith("/auth")) {
                goto("/auth/login");
            }

            // Initialize theme from localStorage if available
            const storedTheme = localStorage.getItem("theme");
            if (storedTheme) {
                theme.set(storedTheme);
            }

            // Subscribe to theme changes and update localStorage
            theme.subscribe((value) => {
                localStorage.setItem("theme", value);
                if (value === "dark") {
                    document.documentElement.classList.add("dark");
                } else {
                    document.documentElement.classList.remove("dark");
                }
            });
        } finally {
            setTimeout(() => {
                isLoading = false;
            }, 500);
        }
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
    <header class="fixed top-0 right-0 m-4 z-50">
        <SettingsModal />
        <Button variant="ghost" size="icon" on:click={toggleTheme}>
            {#if $theme === "dark"}
                <Sun
                    class="h-[1.2rem] w-[1.2rem] rotate-0 scale-100 transition-all dark:-rotate-90 dark:scale-0"
                />
                <Moon
                    class="absolute h-[1.2rem] w-[1.2rem] rotate-90 scale-0 transition-all dark:rotate-0 dark:scale-100"
                />
            {:else}
                <Sun
                    class="h-[1.2rem] w-[1.2rem] rotate-0 scale-100 transition-all"
                />
                <Moon
                    class="absolute h-[1.2rem] w-[1.2rem] rotate-90 scale-0 transition-all"
                />
            {/if}
            <span class="sr-only">Toggle theme</span>
        </Button>
    </header>

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
