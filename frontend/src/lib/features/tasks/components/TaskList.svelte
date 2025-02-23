<script lang="ts">
    import { onMount } from "svelte";
    import { tasksStore } from "../stores/tasks.store";
    import { categories, categoryTextColors } from "../constants";
    import TaskCard from "./TaskCard.svelte";
    import TaskDrawer from "./TaskDrawer.svelte";
    import TaskQuickAdd from "./TaskQuickAdd.svelte";
    import { Search } from "lucide-svelte";
    import { Input } from "$lib/components/ui/input";
    import { Button } from "$lib/components/ui/button";
    import { Plus } from "lucide-svelte";

    let showQuickAdd = false;
    let quickAddCategory = "";
    let searchQuery = "";
    let localTasks = [];

    function handleQuickAdd(category: string) {
        quickAddCategory = category;
        showQuickAdd = true;
    }

    async function handleSearch(e: KeyboardEvent) {
        if (e.key === "Enter") {
            // Server-side search
            await tasksStore.fetchTasks(searchQuery);
        } else {
            // Local search
            if (!searchQuery.trim()) {
                localTasks = $tasksStore.tasks;
            } else {
                const query = searchQuery.toLowerCase();
                localTasks = $tasksStore.tasks.filter(
                    (task) =>
                        task.title.toLowerCase().includes(query) ||
                        task.description.toLowerCase().includes(query),
                );
            }
        }
    }

    $: categoryTasks = categories.map((cat) => ({
        ...cat,
        tasks: localTasks.filter((task) => task.category === cat.value),
    }));

    onMount(() => {
        tasksStore.fetchTasks();
        localTasks = $tasksStore.tasks;
    });

    // Update local tasks when store tasks change
    $: if ($tasksStore.tasks) {
        localTasks = $tasksStore.tasks;
    }
</script>

<div class="container mx-auto p-4">
    <!-- Search Bar -->
    <div class="mb-6 flex items-center gap-4">
        <div class="relative flex-1 max-w-md">
            <Search
                class="absolute left-2 top-2.5 h-4 w-4 text-muted-foreground"
            />
            <Input
                type="text"
                placeholder="Search tasks... (Press Enter for server search)"
                class="pl-8"
                bind:value={searchQuery}
                on:keydown={handleSearch}
                on:input={handleSearch}
            />
        </div>
    </div>

    <!-- Tasks Grid -->
    <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6">
        {#each categoryTasks as category}
            <div class="flex flex-col h-full">
                <div class="flex items-center justify-between mb-4">
                    <h2
                        class="text-lg font-semibold {categoryTextColors[
                            category.value
                        ]}"
                    >
                        {category.label}
                    </h2>
                    <Button
                        variant="ghost"
                        size="sm"
                        class={categoryTextColors[category.value]}
                        on:click={() => handleQuickAdd(category.value)}
                    >
                        <Plus class="h-4 w-4" />
                    </Button>
                </div>

                <div
                    class="flex-1 {category.color} rounded-lg p-4 min-h-[500px]"
                >
                    {#if showQuickAdd && quickAddCategory === category.value}
                        <TaskQuickAdd
                            category={category.value}
                            on:close={() => {
                                showQuickAdd = false;
                                quickAddCategory = "";
                            }}
                        />
                    {/if}

                    {#each category.tasks as task (task.id)}
                        <TaskCard {task} />
                    {/each}
                </div>
            </div>
        {/each}
    </div>
</div>

<TaskDrawer />
