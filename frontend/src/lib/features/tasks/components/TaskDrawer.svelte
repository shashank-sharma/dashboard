<script lang="ts">
    import { format } from "date-fns";
    import * as Drawer from "$lib/components/ui/drawer";
    import { Button } from "$lib/components/ui/button";
    import { Input } from "$lib/components/ui/input";
    import { Textarea } from "$lib/components/ui/textarea";
    import { X, Calendar, Clock, Trash2 } from "lucide-svelte";
    import DateTimePicker from "$lib/components/DateTimePicker.svelte";
    import { tasksStore } from "../stores/tasks.store";
    import { drawerStore } from "../stores/drawer.store";
    import { categories, categoryTextColors } from "../constants";
    import * as AlertDialog from "$lib/components/ui/alert-dialog";

    let taskDueDate: Date | undefined = undefined;
    let showDeleteDialog = false;

    $: if ($drawerStore.task?.due) {
        taskDueDate = new Date($drawerStore.task.due);
    }

    function handleDateSelect(date: Date) {
        taskDueDate = date;
    }

    function handleSave() {
        if (!$drawerStore.task) return;

        tasksStore.updateTask({
            ...$drawerStore.task,
            due: taskDueDate?.toISOString(),
        });
        drawerStore.close();
    }

    function handleDelete() {
        if (!$drawerStore.task) return;
        tasksStore.deleteTask($drawerStore.task.id);
        drawerStore.close();
    }

    function handleDrawerClose() {
        drawerStore.close();
    }
</script>

<Drawer.Root bind:open={$drawerStore.isOpen} onOpenChange={handleDrawerClose}>
    <Drawer.Content class="h-[75%]">
        <div class="mx-auto w-full max-w-2xl h-full relative">
            <Drawer.Close class="absolute right-4 top-4">
                <Button
                    variant="ghost"
                    size="icon"
                    on:click={() => drawerStore.close()}
                >
                    <X class="h-4 w-4" />
                    <span class="sr-only">Close</span>
                </Button>
            </Drawer.Close>

            <Drawer.Header>
                <Drawer.Title>Edit Task</Drawer.Title>
                <Drawer.Description>
                    Make changes to your task here. Click save when you're done.
                </Drawer.Description>
            </Drawer.Header>

            {#if $drawerStore.task}
                <div class="overflow-y-auto px-4 h-[calc(100%-8rem)]">
                    <div class="space-y-4">
                        <!-- Title -->
                        <div class="space-y-2">
                            <label for="title" class="text-sm font-medium"
                                >Title</label
                            >
                            <Input
                                id="title"
                                value={$drawerStore.task.title}
                                on:input={(e) =>
                                    drawerStore.updateTask({
                                        title: e.target.value,
                                    })}
                                placeholder="Task title"
                            />
                        </div>

                        <!-- Description -->
                        <div class="space-y-2">
                            <label for="description" class="text-sm font-medium"
                                >Description</label
                            >
                            <Textarea
                                id="description"
                                value={$drawerStore.task.description}
                                on:input={(e) =>
                                    drawerStore.updateTask({
                                        description: e.target.value,
                                    })}
                                placeholder="Add a more detailed description..."
                                rows="4"
                            />
                        </div>

                        <!-- Category -->
                        <div class="space-y-2">
                            <label class="text-sm font-medium">Category</label>
                            <div class="flex gap-2 flex-wrap">
                                {#each categories as cat}
                                    <Button
                                        variant={$drawerStore.task.category ===
                                        cat.value
                                            ? "default"
                                            : "outline"}
                                        size="sm"
                                        class={$drawerStore.task.category ===
                                        cat.value
                                            ? categoryTextColors[cat.value]
                                            : ""}
                                        on:click={() =>
                                            drawerStore.updateTask({
                                                category: cat.value,
                                            })}
                                    >
                                        {cat.label}
                                    </Button>
                                {/each}
                            </div>
                        </div>

                        <!-- Due Date -->
                        <div class="space-y-2">
                            <label class="text-sm font-medium">Due Date</label>
                            <div class="flex gap-2 items-center">
                                <DateTimePicker
                                    value={taskDueDate}
                                    on:change={(e) =>
                                        handleDateSelect(e.detail)}
                                    placeholder="Select due date and time"
                                />
                                {#if taskDueDate}
                                    <Button
                                        variant="ghost"
                                        size="icon"
                                        on:click={() =>
                                            (taskDueDate = undefined)}
                                    >
                                        <X class="h-4 w-4" />
                                    </Button>
                                {/if}
                            </div>
                        </div>

                        <!-- Metadata -->
                        <div class="space-y-2 pt-4 border-t">
                            <h4
                                class="text-sm font-medium text-muted-foreground"
                            >
                                Task Information
                            </h4>
                            <div class="space-y-2">
                                <div
                                    class="flex items-center gap-2 text-sm text-muted-foreground"
                                >
                                    <Calendar class="h-4 w-4" />
                                    Created: {format(
                                        new Date($drawerStore.task.created),
                                        "PPP",
                                    )}
                                </div>
                                <div
                                    class="flex items-center gap-2 text-sm text-muted-foreground"
                                >
                                    <Clock class="h-4 w-4" />
                                    Last updated: {format(
                                        new Date($drawerStore.task.updated),
                                        "PPP",
                                    )}
                                </div>
                            </div>
                        </div>
                    </div>
                </div>

                <Drawer.Footer>
                    <div class="flex w-full justify-between">
                        <Button
                            variant="destructive"
                            class="gap-2"
                            on:click={() => (showDeleteDialog = true)}
                        >
                            <Trash2 class="h-4 w-4" />
                            Delete
                        </Button>
                        <div class="flex gap-2">
                            <Drawer.Close>
                                <Button
                                    variant="outline"
                                    on:click={() => drawerStore.close()}
                                >
                                    Cancel
                                </Button>
                            </Drawer.Close>
                            <Button on:click={handleSave}>Save changes</Button>
                        </div>
                    </div>
                </Drawer.Footer>
            {/if}
        </div>
    </Drawer.Content>
</Drawer.Root>

<!-- Delete Confirmation Dialog -->
<AlertDialog.Root bind:open={showDeleteDialog}>
    <AlertDialog.Content>
        <AlertDialog.Header>
            <AlertDialog.Title>Delete Task</AlertDialog.Title>
            <AlertDialog.Description>
                Are you sure you want to delete this task? This action cannot be
                undone.
            </AlertDialog.Description>
        </AlertDialog.Header>
        <AlertDialog.Footer>
            <AlertDialog.Cancel on:click={() => (showDeleteDialog = false)}>
                Cancel
            </AlertDialog.Cancel>
            <AlertDialog.Action
                on:click={handleDelete}
                class="bg-destructive text-destructive-foreground"
            >
                Delete
            </AlertDialog.Action>
        </AlertDialog.Footer>
    </AlertDialog.Content>
</AlertDialog.Root>
