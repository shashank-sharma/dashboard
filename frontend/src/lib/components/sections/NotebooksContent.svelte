<script lang="ts">
    import { onMount } from "svelte";
    import { goto } from "$app/navigation";
    import { notebookService } from "$lib/services/notebookService";
    import { Button } from "$lib/components/ui/button";
    import {
        Dialog,
        DialogContent,
        DialogHeader,
        DialogFooter,
        DialogTitle,
        DialogDescription,
    } from "$lib/components/ui/dialog";
    import { Input } from "$lib/components/ui/input";
    import { Label } from "$lib/components/ui/label";
    import {
        Card,
        CardContent,
        CardFooter,
        CardHeader,
        CardTitle,
        CardDescription,
    } from "$lib/components/ui/card";
    import {
        Plus,
        Trash2,
        Search,
        ChevronLeft,
        ChevronRight,
    } from "lucide-svelte";
    import { toast } from "svelte-sonner";

    let notebooks: any[] = [];
    let loading = true;
    let dialogOpen = false;
    let newNotebookName = "";
    let selectedVersion = "python3";
    let searchQuery = "";
    let currentPage = 1;
    let itemsPerPage = 6;
    let totalPages = 0;

    $: filteredNotebooks = notebooks.filter((notebook) =>
        notebook.name.toLowerCase().includes(searchQuery.toLowerCase()),
    );

    $: {
        totalPages = Math.ceil(filteredNotebooks.length / itemsPerPage);
        if (currentPage > totalPages) currentPage = totalPages || 1;
    }

    $: paginatedNotebooks = filteredNotebooks.slice(
        (currentPage - 1) * itemsPerPage,
        currentPage * itemsPerPage,
    );

    async function loadNotebooks() {
        try {
            notebooks = await notebookService.listNotebooks();
        } catch (error) {
            toast.error("Failed to load notebooks");
        } finally {
            loading = false;
        }
    }

    async function handleCreateNotebook() {
        if (!newNotebookName.trim()) return;

        try {
            const notebook = await notebookService.createNotebook({
                name: newNotebookName,
                version: selectedVersion,
                cells: [
                    {
                        id: window.crypto.randomUUID(),
                        content:
                            "# Welcome to your new notebook\n# Press Run to execute code",
                        output: "",
                        type: "code",
                        language: "python",
                    },
                ],
            });

            notebooks = [notebook, ...notebooks];
            dialogOpen = false;
            newNotebookName = "";
            toast.success("Notebook created successfully");
        } catch (error) {
            console.log("Error creating notebook: ", error);
            toast.error("Failed to create notebook");
        }
    }

    async function deleteNotebook(id: string) {
        if (!confirm("Are you sure you want to delete this notebook?")) return;

        try {
            await notebookService.deleteNotebook(id);
            notebooks = notebooks.filter((nb) => nb.id !== id);
            toast.success("Notebook deleted successfully");
        } catch (error) {
            toast.error("Failed to delete notebook");
        }
    }

    onMount(loadNotebooks);
</script>

<div class="space-y-6">
    <!-- Header with Create Button -->
    <div class="flex justify-between items-center">
        <h2 class="text-2xl font-bold">My Notebooks</h2>
        <Dialog bind:open={dialogOpen}>
            <Button
                on:click={() => (dialogOpen = true)}
                class="flex items-center gap-2"
            >
                <Plus class="w-4 h-4" />
                New Notebook
            </Button>
            <DialogContent>
                <DialogHeader>
                    <DialogTitle>Create New Notebook</DialogTitle>
                    <DialogDescription>
                        Create a new notebook to start coding and experimenting.
                    </DialogDescription>
                </DialogHeader>
                <div class="grid gap-4 py-4">
                    <div class="grid gap-2">
                        <Label for="name">Notebook Name</Label>
                        <Input
                            id="name"
                            bind:value={newNotebookName}
                            placeholder="My Awesome Notebook"
                        />
                    </div>
                    <div class="grid gap-2">
                        <Label for="version">Python Version</Label>
                        <select
                            id="version"
                            bind:value={selectedVersion}
                            class="flex h-9 w-full rounded-md border border-input bg-background px-3 py-1 text-sm shadow-sm transition-colors file:border-0 file:bg-transparent file:text-sm file:font-medium placeholder:text-muted-foreground focus-visible:outline-none focus-visible:ring-1 focus-visible:ring-ring disabled:cursor-not-allowed disabled:opacity-50"
                        >
                            <option value="python3">Python 3</option>
                            <option value="python2">Python 2</option>
                        </select>
                    </div>
                </div>
                <DialogFooter>
                    <Button on:click={handleCreateNotebook}
                        >Create Notebook</Button
                    >
                </DialogFooter>
            </DialogContent>
        </Dialog>
    </div>

    <!-- Search Bar -->
    <div class="flex gap-4 items-center">
        <div class="relative flex-1">
            <div class="absolute left-2 top-2.5 h-4 w-4 text-muted-foreground">
                <Search size={16} />
            </div>
            <Input
                type="text"
                placeholder="Search notebooks..."
                bind:value={searchQuery}
                class="pl-8"
            />
        </div>
    </div>

    <!-- Content -->
    {#if loading}
        <div class="text-center py-12">
            <p class="text-muted-foreground">Loading notebooks...</p>
        </div>
    {:else}
        <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
            {#each paginatedNotebooks as notebook (notebook.id)}
                <Card class="hover:shadow-lg transition-shadow">
                    <CardHeader>
                        <CardTitle>{notebook.name}</CardTitle>
                        <CardDescription>
                            Version: {notebook.version}<br />
                            Created: {new Date(
                                notebook.created,
                            ).toLocaleDateString()}
                        </CardDescription>
                    </CardHeader>
                    <CardContent>
                        <p class="text-sm text-muted-foreground">
                            {notebook.cells.length} cells<br />
                            Last modified: {new Date(
                                notebook.updated,
                            ).toLocaleDateString()}
                        </p>
                    </CardContent>
                    <CardFooter class="flex justify-between">
                        <Button
                            variant="outline"
                            on:click={() =>
                                goto(`/dashboard/notebooks/${notebook.id}`)}
                        >
                            Open
                        </Button>
                        <Button
                            variant="destructive"
                            size="icon"
                            on:click={() => deleteNotebook(notebook.id)}
                        >
                            <Trash2 class="w-4 h-4" />
                        </Button>
                    </CardFooter>
                </Card>
            {/each}
        </div>

        <!-- Empty State -->
        {#if filteredNotebooks.length === 0}
            <div class="text-center py-12">
                <p class="text-muted-foreground">
                    {searchQuery
                        ? "No notebooks found matching your search."
                        : "No notebooks yet. Create one to get started!"}
                </p>
            </div>
        {/if}

        <!-- Pagination -->
        {#if totalPages > 1}
            <div class="flex justify-center gap-2 mt-4">
                <Button
                    variant="outline"
                    size="icon"
                    disabled={currentPage === 1}
                    on:click={() => currentPage--}
                >
                    <ChevronLeft class="w-4 h-4" />
                </Button>

                <span class="flex items-center px-4 text-sm">
                    Page {currentPage} of {totalPages}
                </span>

                <Button
                    variant="outline"
                    size="icon"
                    disabled={currentPage === totalPages}
                    on:click={() => currentPage++}
                >
                    <ChevronRight class="w-4 h-4" />
                </Button>
            </div>
        {/if}
    {/if}
</div>
