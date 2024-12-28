<script lang="ts">
    import { Button } from "$lib/components/ui/button";
    import { onMount } from "svelte";
    import { Avatar, AvatarImage, AvatarFallback } from "$lib/components/ui/avatar";
    import { slide } from 'svelte/transition';
    import { tweened } from 'svelte/motion';
    import { cubicOut } from 'svelte/easing';
    
    export let sections: { id: string; label: string; icon: string }[];
    export let activeSection: string;
    export let setActiveSection: (sectionId: string) => void;
    export let username: string = "User"; // Need to pass as prop
    export let avatarUrl: string | null = null; // Need to pass as prop

    let Icons;
    let mounted = false;
    let isCollapsed = false;

    // Width animation
    const width = tweened(256, {
        duration: 200,
        easing: cubicOut
    });

    $: width.set(isCollapsed ? 84 : 256);

    onMount(async () => {
        Icons = await import("lucide-svelte");
        mounted = true;
    });

    function toggleSidebar() {
        isCollapsed = !isCollapsed;
    }
</script>

<aside 
    class="bg-card text-card-foreground border-r flex flex-col h-screen transition-all duration-200 ease-out"
    style:width="{$width}px"
>
    <!-- Dashboard Header -->
    <div class="p-4 border-b">
        <Button 
            variant="ghost" 
            class="w-full justify-start" 
            on:click={toggleSidebar}
        >
            {#if mounted && Icons['LayoutDashboard']}
                <svelte:component
                    this={Icons['LayoutDashboard']}
                    class="h-5 w-5"
                />
            {/if}
            {#if !isCollapsed}
                <span class="ml-2" transition:slide|local>Dashboard</span>
            {/if}
        </Button>
    </div>

    <!-- Navigation -->
    <nav class="p-4 space-y-2 flex-1">
        {#each sections as section}
            <Button
                variant={activeSection === section.id ? "secondary" : "ghost"}
                class="w-full justify-start"
                on:click={() => setActiveSection(section.id)}
            >
                {#if mounted && Icons[section.icon]}
                    <svelte:component
                        this={Icons[section.icon]}
                        class="h-4 w-4"
                    />
                {/if}
                {#if !isCollapsed}
                    <span class="ml-2" transition:slide|local>{section.label}</span>
                {/if}
            </Button>
        {/each}
    </nav>

    <!-- User Profile -->
    <div class="p-4 border-t">
        <div class="flex items-center gap-2">
            <Avatar class="h-8 w-8">
                {#if avatarUrl}
                    <AvatarImage src={avatarUrl} alt={username} />
                {/if}
                <AvatarFallback>
                    {username.slice(0, 2).toUpperCase()}
                </AvatarFallback>
            </Avatar>
            {#if !isCollapsed}
                <span class="text-sm font-medium truncate" transition:slide|local>
                    {username}
                </span>
            {/if}
        </div>
    </div>
</aside>

<style>
    /* Hide scrollbar but keep functionality */
    nav {
        scrollbar-width: none;
        -ms-overflow-style: none;
        overflow-y: auto;
    }
    nav::-webkit-scrollbar {
        display: none;
    }
</style>