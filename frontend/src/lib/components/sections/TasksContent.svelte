<script lang="ts">
    import { onMount } from "svelte";
    import { pb } from "$lib/pocketbase";
    import { toast } from "svelte-sonner";
    import { Input } from "$lib/components/ui/input";
    import { Button } from "$lib/components/ui/button";
    import { Textarea } from "$lib/components/ui/textarea";
    import * as Drawer from "$lib/components/ui/drawer";
    import * as AlertDialog from "$lib/components/ui/alert-dialog";
    import { Search, Trash2, Plus, X, Calendar, Clock } from "lucide-svelte";
    import { fade, fly } from "svelte/transition";
    import { format } from "date-fns";
    import DateTimePicker from "$lib/components/DateTimePicker.svelte";

    // Types
    interface Category {
        value: string;
        label: string;
        color: string;
    }

    interface Task {
        id: string;
        title: string;
        description: string;
        category: string;
        due?: string;
        created: string;
        updated: string;
        user: string;
    }

    // Constants
    const categories: Category[] = [
        {
            value: "focus",
            label: "Focus",
            color: "bg-red-100 dark:bg-red-900/20",
        },
        {
            value: "goals",
            label: "Goals",
            color: "bg-blue-100 dark:bg-blue-900/20",
        },
        {
            value: "fitin",
            label: "Fit In",
            color: "bg-green-100 dark:bg-green-900/20",
        },
        {
            value: "backburner",
            label: "Backburner",
            color: "bg-yellow-100 dark:bg-yellow-900/20",
        },
    ];

    const categoryTextColors: Record<string, string> = {
        focus: "text-red-700 dark:text-red-300",
        goals: "text-blue-700 dark:text-blue-300",
        fitin: "text-green-700 dark:text-green-300",
        backburner: "text-yellow-700 dark:text-yellow-300",
    };

    // State
    let tasks: Task[] = [];
    let searchQuery = "";
    let showDeleteDialog = false;
    let taskToDelete: Task | null = null;
    let selectedTask: Task | null = null;
    let showDrawer = false;
    let taskDueDate: Date | undefined = undefined;
    let showQuickAdd = false;
    let quickAddCategory = "";
    let newTaskTitle = "";

    // Task operations
    async function fetchTasks(): Promise<void> {
        try {
            const filterConditions = [`user = "${pb.authStore.model?.id}"`];
            if (searchQuery) {
                filterConditions.push(
                    `(title ~ "${searchQuery}" || description ~ "${searchQuery}")`,
                );
            }

            const filter = filterConditions.join(" && ");

            const resultList = await pb.collection("tasks").getList(1, 50, {
                sort: "-created",
                filter,
            });

            tasks = resultList.items;
        } catch (error) {
            console.error("Error fetching tasks:", error);
            toast.error("Failed to fetch tasks");
        }
    }

    async function updateTask(task: Task, fullUpdate = false): Promise<void> {
        try {
            const updateData = fullUpdate
                ? {
                      title: task.title,
                      description: task.description,
                      category: task.category,
                      due: taskDueDate?.toISOString(),
                  }
                : {
                      title: task.title,
                      description: task.description,
                      category: task.category,
                  };

            await pb.collection("tasks").update(task.id, updateData);
            await fetchTasks();
            toast.success("Task updated");
        } catch (error) {
            console.error("Error updating task:", error);
            toast.error("Failed to update task");
        }
    }

    async function deleteTask(): Promise<void> {
        if (!taskToDelete) return;

        try {
            await pb.collection("tasks").delete(taskToDelete.id);
            toast.success("Task deleted");
            showDeleteDialog = false;
            taskToDelete = null;
            await fetchTasks();
        } catch (error) {
            console.error("Error deleting task:", error);
            toast.error("Failed to delete task");
        }
    }

    async function quickAddTask(): Promise<void> {
        if (!newTaskTitle.trim() || !quickAddCategory) return;

        try {
            await pb.collection("tasks").create({
                title: newTaskTitle.trim(),
                description: "",
                category: quickAddCategory,
                user: pb.authStore.model?.id,
            });

            newTaskTitle = "";
            showQuickAdd = false;
            quickAddCategory = "";
            await fetchTasks();
            toast.success("Task created");
        } catch (error) {
            console.error("Error creating task:", error);
            toast.error("Failed to create task");
        }
    }

    // Event handlers
    function handleQuickAdd(category: string): void {
        quickAddCategory = category;
        showQuickAdd = true;
    }

    function handleCardClick(task: Task): void {
        selectedTask = { ...task };
        taskDueDate = task.due ? new Date(task.due) : undefined;
        showDrawer = true;
    }

    function handleDateSelect(date: Date) {
        taskDueDate = date;
    }

    // Computed
    $: categoryTasks = categories.map((cat) => ({
        ...cat,
        tasks: tasks.filter((task) => task.category === cat.value),
    }));

    onMount(() => {
        void fetchTasks();
    });
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
                placeholder="Search tasks..."
                class="pl-8"
                bind:value={searchQuery}
                on:input={() => void fetchTasks()}
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
                        <div
                            class="bg-card border rounded-lg p-3 mb-3 shadow-sm"
                            in:fly={{ y: -20, duration: 200 }}
                        >
                            <Input
                                placeholder="Enter task title..."
                                bind:value={newTaskTitle}
                                class="mb-2"
                                on:keydown={(e) =>
                                    e.key === "Enter" && void quickAddTask()}
                            />
                            <div class="flex justify-end gap-2">
                                <Button
                                    variant="ghost"
                                    size="sm"
                                    on:click={() => {
                                        showQuickAdd = false;
                                        newTaskTitle = "";
                                    }}
                                >
                                    Cancel
                                </Button>
                                <Button
                                    size="sm"
                                    on:click={() => void quickAddTask()}
                                >
                                    Add
                                </Button>
                            </div>
                        </div>
                    {/if}

                    {#each category.tasks as task (task.id)}
                        <div
                            class="bg-card border rounded-lg p-3 mb-3 shadow-sm hover:shadow-md transition-shadow cursor-pointer"
                            in:fade
                            on:click={() => handleCardClick(task)}
                            on:keydown={(e) =>
                                e.key === "Enter" && handleCardClick(task)}
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
                                    on:click={(e) => {
                                        e.stopPropagation();
                                        taskToDelete = task;
                                        showDeleteDialog = true;
                                    }}
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
                    {/each}
                </div>
            </div>
        {/each}
    </div>
</div>

<!-- Task Details Drawer -->
<Drawer.Root bind:open={showDrawer}>
    <Drawer.Content class="h-[75%]">
        <div class="mx-auto w-full max-w-2xl h-full relative">
            <Drawer.Close class="absolute right-4 top-4">
                <Button variant="ghost" size="icon">
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

            {#if selectedTask}
                <div class="overflow-y-auto px-4 h-[calc(100%-8rem)]">
                    <div class="space-y-4">
                        <!-- Title -->
                        <div class="space-y-2">
                            <label for="title" class="text-sm font-medium"
                                >Title</label
                            >
                            <Input
                                id="title"
                                bind:value={selectedTask.title}
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
                                bind:value={selectedTask.description}
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
                                        variant={selectedTask.category ===
                                        cat.value
                                            ? "default"
                                            : "outline"}
                                        size="sm"
                                        class={selectedTask.category ===
                                        cat.value
                                            ? categoryTextColors[cat.value]
                                            : ""}
                                        on:click={() =>
                                            (selectedTask.category = cat.value)}
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
                                        new Date(selectedTask.created),
                                        "PPP",
                                    )}
                                </div>
                                <div
                                    class="flex items-center gap-2 text-sm text-muted-foreground"
                                >
                                    <Clock class="h-4 w-4" />
                                    Last updated: {format(
                                        new Date(selectedTask.updated),
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
                            on:click={() => {
                                taskToDelete = selectedTask;
                                showDrawer = false;
                                showDeleteDialog = true;
                            }}
                        >
                            <Trash2 class="h-4 w-4" />
                            Delete
                        </Button>
                        <div class="flex gap-2">
                            <Drawer.Close>
                                <Button variant="outline">Cancel</Button>
                            </Drawer.Close>
                            <Button
                                on:click={() => {
                                    void updateTask(selectedTask, true);
                                    showDrawer = false;
                                }}
                            >
                                Save changes
                            </Button>
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
            <AlertDialog.Cancel>Cancel</AlertDialog.Cancel>
            <AlertDialog.Action
                on:click={() => void deleteTask()}
                class="bg-destructive text-destructive-foreground"
            >
                Delete
            </AlertDialog.Action>
        </AlertDialog.Footer>
    </AlertDialog.Content>
</AlertDialog.Root>
