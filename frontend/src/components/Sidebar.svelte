<script lang="ts">
    import { createEventDispatcher } from "svelte";
    import { tweened } from "svelte/motion";
    import { cubicOut } from "svelte/easing";
    import { fade, fly } from "svelte/transition";
    import { page } from "$app/stores";
    import {
        Icon,
        ArrowLeft,
        ArrowRight,
        DocumentText,
        CurrencyDollar,
        ChartBar,
        Server,
    } from "svelte-hero-icons";

    export let isSidebarOpen;

    const dispatch = createEventDispatcher();

    interface Section {
        name: string;
        icon: any;
        path: string;
    }

    const sections: Section[] = [
        { name: "Notes", icon: DocumentText, path: "/dashboard/notes" },
        { name: "Tasks", icon: DocumentText, path: "/dashboard/tasks" },
        { name: "Money", icon: CurrencyDollar, path: "/dashboard/money" },
        { name: "Monitor", icon: ChartBar, path: "/dashboard/monitor" },
        { name: "Servers", icon: Server, path: "/dashboard/servers" },
    ];

    const sidebarWidth = tweened(60, {
        duration: 300,
        easing: cubicOut,
    });

    $: {
        $sidebarWidth = $isSidebarOpen ? 240 : 60;
    }

    function toggleSidebar() {
        isSidebarOpen.update((value) => !value);
    }
</script>

<aside class="sidebar" style="width: {$sidebarWidth}px">
    <div class="sidebar-header">
        <button
            class="toggle-btn"
            on:click={toggleSidebar}
            aria-label={$isSidebarOpen ? "Collapse sidebar" : "Expand sidebar"}
        >
            <Icon src={$isSidebarOpen ? ArrowLeft : ArrowRight} size="24" />
        </button>
        {#if $isSidebarOpen}
            <h2 in:fade={{ duration: 200, delay: 100 }}>Dashboard</h2>
        {/if}
    </div>
    <nav>
        {#each sections as section (section.name)}
            <a
                href={section.path}
                class="section-btn"
                class:active={$page.url.pathname === section.path}
                aria-current={$page.url.pathname === section.path
                    ? "page"
                    : undefined}
                title={section.name}
            >
                <Icon src={section.icon} size="24" />
                {#if $isSidebarOpen}
                    <span
                        class="section-name"
                        in:fly={{ x: -20, duration: 200, delay: 100 }}
                    >
                        {section.name}
                    </span>
                {/if}
            </a>
        {/each}
    </nav>
</aside>

<style>
    :root {
        --sidebar-bg: white;
        --sidebar-text: #333;
        --sidebar-hover-bg: rgba(0, 0, 0, 0.05);
        --sidebar-active-color: #3498db;
        --sidebar-border-color: #e0e0e0;
        --sidebar-transition: 0.3s cubic-bezier(0.4, 0, 0.2, 1);
    }

    .sidebar {
        background-color: var(--sidebar-bg);
        color: var(--sidebar-text);
        height: 100%;
        display: flex;
        flex-direction: column;
        overflow: hidden;
        box-shadow: 0 0 10px rgba(0, 0, 0, 0.1);
        transition: width var(--sidebar-transition);
        border-radius: 0 16px 16px 0;
    }

    .sidebar-header {
        display: flex;
        align-items: center;
        padding: 1rem;
        border-bottom: 1px solid var(--sidebar-border-color);
    }

    .sidebar-header h2 {
        margin: 0;
        font-size: 1.2rem;
        font-weight: 600;
        white-space: nowrap;
        overflow: hidden;
        text-overflow: ellipsis;
    }

    .toggle-btn {
        background: none;
        border: none;
        color: var(--sidebar-text);
        cursor: pointer;
        padding: 0.5rem;
        margin-right: 0.5rem;
        border-radius: 50%;
        transition: background-color 0.2s ease;
    }

    .toggle-btn:hover {
        background-color: var(--sidebar-hover-bg);
    }

    nav {
        display: flex;
        flex-direction: column;
        padding: 1rem 0;
    }

    .section-btn {
        display: flex;
        align-items: center;
        padding: 0.75rem 1rem;
        background: none;
        border-right: 3px solid transparent;
        color: var(--sidebar-text);
        cursor: pointer;
        transition: all 0.2s ease;
        margin: 0.25rem 0;
        text-decoration: none;
    }

    .section-btn:hover,
    .section-btn.active {
        border-right-color: var(--sidebar-active-color);
        color: var(--sidebar-active-color);
    }

    .section-btn.active {
        font-weight: bold;
    }

    .section-name {
        margin-left: 1rem;
        font-size: 0.9rem;
        white-space: nowrap;
    }

    @media (max-width: 768px) {
        .sidebar {
            border-radius: 0;
        }
    }
</style>
