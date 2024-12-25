<script lang="ts">
    import { Button } from "$lib/components/ui/button";
    import { onMount } from "svelte";

    export let sections: { id: string; label: string; icon: string }[];
    export let activeSection: string;
    export let setActiveSection: (sectionId: string) => void;

    let Icons;
    let mounted = false;

    onMount(async () => {
        Icons = await import("lucide-svelte");
        mounted = true;
    });
</script>

<aside class="w-64 bg-card text-card-foreground border-r">
    <nav class="p-4 space-y-2">
        {#each sections as section}
            <Button
                variant={activeSection === section.id ? "secondary" : "ghost"}
                class="w-full justify-start"
                on:click={() => setActiveSection(section.id)}
            >
                {#if mounted && Icons[section.icon]}
                    <svelte:component
                        this={Icons[section.icon]}
                        class="mr-2 h-4 w-4"
                    />
                {/if}
                {section.label}
            </Button>
        {/each}
    </nav>
</aside>
