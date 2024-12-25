<script lang="ts">
    import { onMount } from "svelte";
    import { goto } from "$app/navigation";
    import { page } from "$app/stores";
    import Sidebar from "./Sidebar.svelte";
    import Content from "./Content.svelte";

    let activeSection = "dashboard";

    const sections = [
        { id: "dashboard", label: "Dashboard", icon: "LayoutDashboard" },
        { id: "servers", label: "Servers", icon: "Server" },
        { id: "notebooks", label: "Notebooks", icon: "Book" },
        { id: "posts", label: "Posts", icon: "FileText" },
        { id: "colors", label: "Colors", icon: "Palette" },
        { id: "inventory", label: "Inventory", icon: "Package" },
        { id: "newsletter", label: "Newsletter", icon: "Mail" },
        { id: "bookmarks", label: "Bookmarks", icon: "Bookmark" },
        { id: "expenses", label: "Expenses", icon: "DollarSign" },
        { id: "tasks", label: "Tasks", icon: "CheckSquare" },
    ];

    function setActiveSection(sectionId: string) {
        activeSection = sectionId;
        goto(`/dashboard/${sectionId}`);
    }

    $: {
        const path = $page.url.pathname;
        const section = path.split("/").pop();
        if (sections.some((s) => s.id === section)) {
            activeSection = section;
        }
    }

    onMount(() => {
        const path = $page.url.pathname;
        const section = path.split("/").pop();
        if (section && section !== "dashboard") {
            setActiveSection(section);
        }
    });
</script>

<div class="flex h-screen bg-background">
    <Sidebar {sections} {activeSection} {setActiveSection} />
    <Content {activeSection} />
</div>
