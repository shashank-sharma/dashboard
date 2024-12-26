<script lang="ts">
    import { page } from "$app/stores";
    import { goto } from "$app/navigation";
    import Sidebar from "$lib/components/Sidebar.svelte";
    import SectionLoading from "$lib/components/SectionLoading.svelte";
    import { fade } from "svelte/transition";

    let activeSection = "dashboard";
    let isLoading = false;

    const sections = [
        { id: "", label: "Dashboard", icon: "LayoutDashboard" },
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
    ];

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
</div>

<style>
    main {
        height: calc(100vh - 2rem);
    }
</style>
