<script lang="ts">
    import { Button } from "$lib/components/ui/button";
    import { Plus } from "lucide-svelte";
    import { fade } from "svelte/transition";

    export let title: string;
    export let onAdd: () => void = () => {};
    export let showAddButton = true;
</script>

<div class="space-y-4" in:fade>
    <div class="flex justify-between items-center">
        <h2 class="text-3xl font-bold tracking-tight">{title}</h2>
        {#if showAddButton}
            <Button on:click={onAdd}>
                <Plus class="h-4 w-4 mr-2" />
                Add New
            </Button>
        {/if}
    </div>
    <slot />
</div>
