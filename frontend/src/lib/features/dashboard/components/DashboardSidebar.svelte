<script lang="ts">
    import { page } from "$app/stores";
    import { goto } from "$app/navigation";
    import { Button } from "$lib/components/ui/button";
    import { cn } from "$lib/utils";
    import { auth } from "$lib/stores/auth.store";
    import { dashboardStore } from "../stores/dashboard.store";
    import type { DashboardSection } from "../types";
    import { ChevronLeft, ChevronRight } from "lucide-svelte";
    import { slide } from "svelte/transition";
    import { tweened } from "svelte/motion";
    import { cubicOut } from "svelte/easing";

    export let sections: DashboardSection[];

    let isCollapsed = false;
    const width = tweened(256, {
        duration: 200,
        easing: cubicOut,
    });

    $: width.set(isCollapsed ? 84 : 256);

    function toggleSidebar() {
        isCollapsed = !isCollapsed;
    }

    function navigateToSection(section: DashboardSection) {
        dashboardStore.setActiveSection(section.id);
        goto(section.path);
    }

    $: currentPath = $page.url.pathname;
</script>

<aside
    class="bg-card border-r flex flex-col h-screen transition-all duration-200 ease-out"
    style:width="{$width}px"
>
    <div class="p-4 flex items-center justify-between">
        {#if !isCollapsed}
            <h1 class="text-xl font-bold">Dashboard</h1>
        {/if}
        <Button variant="ghost" size="icon" on:click={toggleSidebar}>
            <svelte:component
                this={isCollapsed ? ChevronRight : ChevronLeft}
                class="h-4 w-4"
            />
        </Button>
    </div>

    <nav class="flex-1 px-2 py-4 space-y-1">
        {#each sections as section}
            <Button
                variant={currentPath === section.path ? "secondary" : "ghost"}
                class={cn(
                    "w-full justify-start",
                    isCollapsed && "justify-center px-2",
                )}
                on:click={() => navigateToSection(section)}
            >
                {#if section.icon}
                    <svelte:component
                        this={section.icon}
                        class="h-4 w-4 mr-2"
                    />
                {/if}
                {#if !isCollapsed}
                    <span>{section.label}</span>
                {/if}
            </Button>
        {/each}
    </nav>

    <div class="p-4 border-t">
        {#if !isCollapsed}
            <div class="flex items-center space-x-3">
                <div class="flex-1 min-w-0">
                    <p class="text-sm font-medium truncate">
                        {$auth.user?.username || "Guest"}
                    </p>
                    <p class="text-xs text-muted-foreground truncate">
                        {$auth.user?.email || ""}
                    </p>
                </div>
            </div>
        {/if}
    </div>
</aside>
