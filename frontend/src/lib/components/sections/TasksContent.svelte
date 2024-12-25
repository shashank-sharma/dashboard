<script lang="ts">
    import { onMount } from "svelte";
    import { pb } from "$lib/pocketbase";
    import { toast } from "svelte-sonner";
    import { Input } from "$lib/components/ui/input";
    import { Button } from "$lib/components/ui/button";
    import { Textarea } from "$lib/components/ui/textarea";
    import {
        Select,
        SelectContent,
        SelectItem,
        SelectTrigger,
        SelectValue,
    } from "$lib/components/ui/select";
    import {
        Card,
        CardContent,
        CardDescription,
        CardHeader,
        CardTitle,
    } from "$lib/components/ui/card";
    import { Badge } from "$lib/components/ui/badge";
    import * as Dialog from "$lib/components/ui/dialog";
    import * as AlertDialog from "$lib/components/ui/alert-dialog";
    import { Search, Trash2, Edit, Plus } from "lucide-svelte";

    let tasks = [];
    let searchQuery = "";
    let selectedCategory = "all";
    let totalPages = 1;
    let currentPage = 1;
    const perPage = 9;

    let showCreateDialog = false;
    let showViewDialog = false;
    let showDeleteDialog = false;
    let selectedTask = null;
    let isEditing = false;

    let formTitle = "";
    let formDescription = "";
    let formCategory = "focus";

    const categories = [
        { value: "focus", label: "Focus" },
        { value: "goals", label: "Goals" },
        { value: "fitin", label: "Fit In" },
        { value: "backburner", label: "Backburner" },
    ];

    const categoryColors = {
        focus: "bg-red-500",
        goals: "bg-blue-500",
        fitin: "bg-green-500",
        backburner: "bg-yellow-500",
    };

    async function fetchTasks() {
        try {
            const filterConditions = [];
            if (selectedCategory !== "all") {
                filterConditions.push(`category = "${selectedCategory}"`);
            }
            if (searchQuery) {
                filterConditions.push(
                    `(title ~ "${searchQuery}" || description ~ "${searchQuery}")`,
                );
            }

            filterConditions.push(`user = "${pb.authStore.model.id}"`);

            const filter = filterConditions.join(" && ");

            const resultList = await pb
                .collection("tasks")
                .getList(currentPage, perPage, {
                    sort: "-created",
                    filter,
                });

            tasks = resultList.items;
            totalPages = Math.ceil(resultList.totalItems / perPage);
        } catch (error) {
            console.error("Error fetching tasks:", error);
            toast.error("Failed to fetch tasks.");
        }
    }

    function handleTaskClick(task) {
        selectedTask = task;
        formTitle = task.title;
        formDescription = task.description;
        formCategory = task.category;
        isEditing = false;
        showViewDialog = true;
    }

    async function createTask() {
        try {
            if (!formTitle.trim()) {
                toast.error("Title is required");
                return;
            }

            await pb.collection("tasks").create({
                title: formTitle.trim(),
                description: formDescription.trim(),
                category: formCategory,
                user: pb.authStore.model.id,
            });

            toast.success("Task created successfully!");
            resetForm();
            showCreateDialog = false;
            fetchTasks();
        } catch (error) {
            console.error("Error creating task:", error);
            toast.error("Failed to create task.");
        }
    }

    async function updateTask() {
        try {
            if (!formTitle.trim()) {
                toast.error("Title is required");
                return;
            }

            await pb.collection("tasks").update(selectedTask.id, {
                title: formTitle.trim(),
                description: formDescription.trim(),
                category: formCategory,
                user: pb.authStore.model.id,
            });

            toast.success("Task updated successfully!");
            resetForm();
            showViewDialog = false;
            fetchTasks();
        } catch (error) {
            console.error("Error updating task:", error);
            toast.error("Failed to update task.");
        }
    }

    async function deleteTask() {
        try {
            await pb.collection("tasks").delete(selectedTask.id);
            toast.success("Task deleted successfully!");
            showDeleteDialog = false;
            showViewDialog = false;
            resetForm();
            fetchTasks();
        } catch (error) {
            console.error("Error deleting task:", error);
            toast.error("Failed to delete task.");
        }
    }

    function resetForm() {
        formTitle = "";
        formDescription = "";
        formCategory = "focus";
        isEditing = false;
        selectedTask = null;
    }

    function handleCategoryChange(value) {
        selectedCategory = value;
        currentPage = 1;
        fetchTasks();
    }

    function handleCloseCreateDialog() {
        showCreateDialog = false;
        resetForm();
    }

    function handleCloseViewDialog() {
        showViewDialog = false;
        isEditing = false;
        resetForm();
    }

    onMount(() => {
        fetchTasks();
    });
</script>

<div class="container mx-auto p-6">
    <div class="flex flex-col gap-6">
        <div
            class="flex flex-col sm:flex-row items-center justify-between gap-4"
        >
            <div class="relative w-full max-w-sm">
                <Search
                    class="absolute left-2 top-2.5 h-4 w-4 text-muted-foreground"
                />
                <Input
                    type="text"
                    placeholder="Search tasks..."
                    class="pl-8"
                    value={searchQuery}
                    on:input={(e) => {
                        searchQuery = e.target.value;
                        fetchTasks();
                    }}
                />
            </div>
            <div class="flex gap-4">
                <Select value={selectedCategory}>
                    <SelectTrigger class="w-[180px]">
                        <SelectValue placeholder="Select category" />
                    </SelectTrigger>
                    <SelectContent>
                        <SelectItem
                            value="all"
                            on:click={() => handleCategoryChange("all")}
                        >
                            All Categories
                        </SelectItem>
                        {#each categories as category}
                            <SelectItem
                                value={category.value}
                                on:click={() =>
                                    handleCategoryChange(category.value)}
                            >
                                {category.label}
                            </SelectItem>
                        {/each}
                    </SelectContent>
                </Select>

                <Button
                    on:click={() => {
                        resetForm();
                        showCreateDialog = true;
                    }}
                >
                    <Plus class="mr-2 h-4 w-4" />
                    New Task
                </Button>
            </div>
        </div>

        <div class="grid gap-6 md:grid-cols-2 lg:grid-cols-3">
            {#each tasks as task (task.id)}
                <div
                    role="button"
                    tabindex="0"
                    on:click={() => handleTaskClick(task)}
                    on:keydown={(e) =>
                        e.key === "Enter" && handleTaskClick(task)}
                >
                    <Card
                        class="cursor-pointer hover:shadow-lg transition-shadow h-full"
                    >
                        <CardHeader>
                            <div class="flex items-center justify-between">
                                <CardTitle class="text-xl"
                                    >{task.title}</CardTitle
                                >
                                <Badge class={categoryColors[task.category]}>
                                    {categories.find(
                                        (c) => c.value === task.category,
                                    )?.label}
                                </Badge>
                            </div>
                            <CardDescription
                                class="text-sm text-muted-foreground"
                            >
                                Created on {new Date(
                                    task.created,
                                ).toLocaleDateString()}
                            </CardDescription>
                        </CardHeader>
                        <CardContent>
                            <p class="text-sm line-clamp-3">
                                {task.description}
                            </p>
                        </CardContent>
                    </Card>
                </div>
            {/each}
        </div>

        {#if tasks.length === 0}
            <div class="text-center py-12">
                <p class="text-muted-foreground">No tasks found.</p>
            </div>
        {/if}

        {#if totalPages > 1}
            <div class="flex justify-center gap-2 mt-6">
                {#each Array(totalPages) as _, i}
                    <Button
                        variant={currentPage === i + 1 ? "default" : "outline"}
                        size="sm"
                        on:click={() => {
                            currentPage = i + 1;
                            fetchTasks();
                        }}
                    >
                        {i + 1}
                    </Button>
                {/each}
            </div>
        {/if}
    </div>
</div>

<!-- Create Task Dialog -->
<Dialog.Root
    open={showCreateDialog}
    onOpenChange={(open) => {
        if (!open) handleCloseCreateDialog();
    }}
>
    <Dialog.Content class="sm:max-w-[625px]">
        <Dialog.Header>
            <Dialog.Title>Create New Task</Dialog.Title>
            <Dialog.Description>
                Fill in the details below to create a new task.
            </Dialog.Description>
        </Dialog.Header>
        <div class="grid gap-4 py-4">
            <Input placeholder="Task title" bind:value={formTitle} />
            <Textarea
                placeholder="Task description"
                bind:value={formDescription}
            />
            <Select value={formCategory}>
                <SelectTrigger>
                    <SelectValue placeholder="Select category" />
                </SelectTrigger>
                <SelectContent>
                    {#each categories as category}
                        <SelectItem
                            value={category.value}
                            on:click={() => (formCategory = category.value)}
                        >
                            {category.label}
                        </SelectItem>
                    {/each}
                </SelectContent>
            </Select>
        </div>
        <Dialog.Footer>
            <Button variant="outline" on:click={handleCloseCreateDialog}>
                Cancel
            </Button>
            <Button on:click={createTask}>Create Task</Button>
        </Dialog.Footer>
    </Dialog.Content>
</Dialog.Root>

<!-- View/Edit Task Dialog -->
<Dialog.Root
    open={showViewDialog}
    onOpenChange={(open) => {
        if (!open) handleCloseViewDialog();
    }}
>
    <Dialog.Content class="sm:max-w-[625px]">
        <Dialog.Header>
            <Dialog.Title>{isEditing ? "Edit Task" : "View Task"}</Dialog.Title>
        </Dialog.Header>
        <div class="grid gap-4 py-4">
            <Input
                placeholder="Task title"
                bind:value={formTitle}
                disabled={!isEditing}
            />
            <Textarea
                placeholder="Task description"
                bind:value={formDescription}
                disabled={!isEditing}
            />
            <Select value={formCategory}>
                <SelectTrigger disabled={!isEditing}>
                    <SelectValue placeholder="Select category" />
                </SelectTrigger>
                <SelectContent>
                    {#each categories as category}
                        <SelectItem
                            value={category.value}
                            on:click={() => (formCategory = category.value)}
                        >
                            {category.label}
                        </SelectItem>
                    {/each}
                </SelectContent>
            </Select>
        </div>
        <Dialog.Footer class="gap-2">
            <Button variant="outline" on:click={handleCloseViewDialog}>
                Close
            </Button>
            {#if isEditing}
                <Button on:click={updateTask}>Save Changes</Button>
            {:else}
                <Button
                    variant="destructive"
                    on:click={() => (showDeleteDialog = true)}
                >
                    <Trash2 class="mr-2 h-4 w-4" />
                    Delete
                </Button>
                <Button on:click={() => (isEditing = true)}>
                    <Edit class="mr-2 h-4 w-4" />
                    Edit
                </Button>
            {/if}
        </Dialog.Footer>
    </Dialog.Content>
</Dialog.Root>

<!-- Delete Confirmation Dialog -->
<AlertDialog.Root
    open={showDeleteDialog}
    onOpenChange={(open) => {
        if (!open) showDeleteDialog = false;
    }}
>
    <AlertDialog.Content>
        <AlertDialog.Header>
            <AlertDialog.Title>Are you sure?</AlertDialog.Title>
            <AlertDialog.Description>
                This action cannot be undone. This will permanently delete the
                task.
            </AlertDialog.Description>
        </AlertDialog.Header>
        <AlertDialog.Footer>
            <AlertDialog.Cancel on:click={() => (showDeleteDialog = false)}>
                Cancel
            </AlertDialog.Cancel>
            <AlertDialog.Action
                on:click={deleteTask}
                class="bg-destructive text-destructive-foreground"
            >
                Delete
            </AlertDialog.Action>
        </AlertDialog.Footer>
    </AlertDialog.Content>
</AlertDialog.Root>
