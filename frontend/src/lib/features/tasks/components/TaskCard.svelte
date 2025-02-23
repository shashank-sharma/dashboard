<script lang="ts">
    import { fade } from "svelte/transition";
    import { Button } from "$lib/components/ui/button";
    import { Trash2 } from "lucide-svelte";
    import { tasksStore } from "../stores/tasks.store";
    import { drawerStore } from "../stores/drawer.store";
    import type { Task } from "../types";

    export let task: Task;

    function handleClick() {
        if (!$tasksStore.selectedTask) {
            drawerStore.open(task);
        }
    }

    function handleDelete(e: MouseEvent) {
        e.stopPropagation();
        tasksStore.setTaskToDelete(task);
    }
</script>

<div
    class="bg-card border rounded-lg p-3 mb-3 shadow-sm hover:shadow-md transition-shadow cursor-pointer"
    in:fade
    on:click={handleClick}
    on:keydown={(e) => e.key === "Enter" && handleClick()}
    role="button"
    tabindex="0"
>
    <div class="flex justify-between items-start gap-2">
        <div class="flex-1">
            {task.title}
        </div>

        <Button
            variant="ghost"
            size="icon"
            class="h-8 w-8 text-destructive hover:text-destructive"
            on:click={handleDelete}
        >
            <Trash2 class="h-4 w-4" />
        </Button>
    </div>

    {#if task.description}
        <p class="text-sm text-muted-foreground mt-2">
            {task.description}
        </p>
    {/if}
</div>
