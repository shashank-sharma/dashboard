<script lang="ts">
    import { page } from "$app/stores";
    import { goto } from "$app/navigation";
    import Sidebar from "$lib/components/Sidebar.svelte";
    import SectionLoading from "$lib/components/SectionLoading.svelte";
    import { fade } from "svelte/transition";
    import CommandPalette from "$lib/components/CommandPalette.svelte";
    import { Button } from "$lib/components/ui/button";
    import { Sun, Moon } from "lucide-svelte";
    import SettingsModal from "$lib/components/SettingsModal.svelte";
    import { writable } from "svelte/store";
    import NotificationCenter from "$lib/components/NotificationCenter.svelte";
    import { onMount, setContext } from "svelte";
    import { pb, refreshToken } from "$lib/pocketbase";

    let activeSection = "dashboard";
    let isLoading = false;
    let theme = writable("light");

    function toggleTheme() {
        theme.update((t) => (t === "dark" ? "light" : "dark"));
    }

    const sections = [
        { id: "", label: "Dashboard", icon: "LayoutDashboard" },
        { id: "mails", label: "Mails", icon: "Inbox" },
        { id: "servers", label: "Servers", icon: "Server" },
        { id: "notebooks", label: "Notebooks", icon: "Book" },
        { id: "posts", label: "Posts", icon: "FileText" },
        { id: "colors", label: "Colors", icon: "Palette" },
        { id: "inventory", label: "Inventory", icon: "Package" },
        { id: "newsletter", label: "Newsletter", icon: "Mail" },
        { id: "bookmarks", label: "Bookmarks", icon: "Bookmark" },
        { id: "expenses", label: "Expenses", icon: "DollarSign" },
        { id: "tasks", label: "Tasks", icon: "CheckSquare" },
        { id: "tokens", label: "Token", icon: "KeySquare" },
        { id: "chronicles", label: "Chronicles", icon: "MountainSnow" },
    ];

    setContext("theme", {
        theme: theme,
        toggleTheme: toggleTheme,
    });

    onMount(async () => {
        try {
            if (pb.authStore.isValid) {
                await refreshToken();
            } else if (!window.location.pathname.startsWith("/auth")) {
                pb.authStore.clear();
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
        }
    });

    async function setActiveSection(sectionId: string) {
        try {
            isLoading = true;
            activeSection = sectionId;
            await goto(`/dashboard/${sectionId}`);
        } finally {
            // Add a small delay to make the loading state visible
            setTimeout(() => {
                isLoading = false;
            }, 300);
        }
    }

    $: {
        const path = $page.url.pathname;
        const section = path.split("/").pop();
        if (sections.some((s) => s.id === section)) {
            activeSection = section;
        }
    }
</script>

<div class="flex h-screen bg-background">
    <Sidebar {sections} {activeSection} {setActiveSection} />
    <main class="flex-1 p-6 overflow-auto relative">
        <h1 class="text-3xl font-bold mb-6 capitalize">
            {activeSection || "Dashboard"}
        </h1>

        {#if isLoading}
            <SectionLoading />
        {:else}
            <div
                in:fade={{ duration: 150, delay: 150 }}
                out:fade={{ duration: 150 }}
            >
                <slot />
            </div>
        {/if}
    </main>
    <CommandPalette />
</div>

<header class="fixed top-0 right-0 m-4 z-2">
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
    <NotificationCenter />
</header>

<style>
    main {
        height: calc(100vh - 2rem);
    }
</style>
