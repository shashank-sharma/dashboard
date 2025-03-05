<script lang="ts">
    import { page } from "$app/stores";
    import { goto } from "$app/navigation";
    import { Button } from "$lib/components/ui/button";
    import { cn } from "$lib/utils";
    import { auth } from "$lib/stores/auth.store";
    import { dashboardStore } from "../stores/dashboard.store";
    import type { DashboardSection } from "../types";
    import { ChevronLeft, ChevronRight, ChevronDown } from "lucide-svelte";
    import { slide } from "svelte/transition";
    import { tweened } from "svelte/motion";
    import { cubicOut } from "svelte/easing";
    import { writable } from "svelte/store";
    import { onMount, afterUpdate } from "svelte";

    export let sections: DashboardSection[];

    let isCollapsed = false;
    let prevPath: string = "";
    const width = tweened(256, {
        duration: 200,
        easing: cubicOut,
    });

    // Store for tracking which sections are expanded
    const expandedSections = writable<Record<string, boolean>>({});

    $: width.set(isCollapsed ? 84 : 256);
    $: currentPath = $page.url.pathname;

    // Reactive statement to update active sections when path changes
    $: if (currentPath !== prevPath) {
        prevPath = currentPath;
        updateActiveSection(currentPath);
    }

    function toggleSidebar() {
        isCollapsed = !isCollapsed;
    }

    function toggleSection(sectionId: string, event?: Event) {
        if (event) {
            event.stopPropagation();
        }
        expandedSections.update((state) => ({
            ...state,
            [sectionId]: !state[sectionId],
        }));
    }

    function navigateToSection(section: DashboardSection) {
        console.log("navigating to section", section);
        console.log("currentPath", currentPath);
        console.log("section.path", section.path);
        dashboardStore.setActiveSection(section.id);
        goto(section.path);
    }

    // Check if the given path is or contains the active path
    function isActive(path: string): boolean {
        return currentPath === path || currentPath.startsWith(`${path}/`);
    }

    // Update the active section based on the current path
    function updateActiveSection(path: string) {
        // Find the active section and set it in the store
        let activeSection = "";
        for (const section of sections) {
            if (section.path === path) {
                activeSection = section.id;
                break;
            }

            if (section.children) {
                for (const child of section.children) {
                    if (child.path === path) {
                        activeSection = child.id;
                        break;
                    }
                }
                if (activeSection) break;
            }
        }

        if (activeSection) {
            dashboardStore.setActiveSection(activeSection);
        }

        // Update expanded sections
        sections.forEach((section) => {
            if (section.collapsible && section.children) {
                const hasActiveChild = section.children.some((child) =>
                    isActive(child.path),
                );
                if (hasActiveChild || section.path === path) {
                    expandedSections.update((state) => ({
                        ...state,
                        [section.id]: true,
                    }));
                }
            }
        });
    }

    onMount(() => {
        // Initialize on component mount
        updateActiveSection(currentPath);
    });
</script>

<aside
    class="bg-card border-r flex flex-col h-screen transition-all duration-200 ease-out"
    style:width="{$width}px"
>
    <div class="p-4 flex items-center justify-between">
        {#if !isCollapsed}
            <h1 class="text-xl font-bold">Nen Space</h1>
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
            <div>
                <Button
                    variant={currentPath === section.path
                        ? "secondary"
                        : "ghost"}
                    class={cn(
                        "w-full justify-start",
                        isCollapsed && "justify-center px-2",
                        isCollapsed && section.collapsible && "relative",
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
                        <span class="flex-1 text-left">{section.label}</span>
                        {#if section.collapsible}
                            <div
                                class="cursor-pointer"
                                on:click={(e) => toggleSection(section.id, e)}
                            >
                                <ChevronDown
                                    class={cn(
                                        "h-4 w-4 transition-transform",
                                        $expandedSections[section.id]
                                            ? "rotate-180"
                                            : "",
                                    )}
                                />
                            </div>
                        {/if}
                    {/if}
                </Button>

                {#if section.collapsible && section.children && $expandedSections[section.id] && !isCollapsed}
                    <div
                        class="ml-6 mt-1 space-y-1"
                        transition:slide={{ duration: 150 }}
                    >
                        {#each section.children as child}
                            <Button
                                variant={currentPath === child.path
                                    ? "secondary"
                                    : "ghost"}
                                class="w-full justify-start"
                                on:click={() => {
                                    dashboardStore.setActiveSection(child.id);
                                    navigateToSection(child);
                                    goto(child.path);
                                }}
                            >
                                {#if child.icon}
                                    <svelte:component
                                        this={child.icon}
                                        class="h-4 w-4 mr-2"
                                    />
                                {/if}
                                <span>{child.label}</span>
                            </Button>
                        {/each}
                    </div>
                {/if}
            </div>
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
