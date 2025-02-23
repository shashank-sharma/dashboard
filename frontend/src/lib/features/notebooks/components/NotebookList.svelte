<script lang="ts">
    import { onMount } from "svelte";
    import { notebooksStore } from "../stores/notebooks.store";
    import { Button } from "$lib/components/ui/button";
    import { Input } from "$lib/components/ui/input";
    import { Search, Plus, ChevronLeft, ChevronRight } from "lucide-svelte";
    import NotebookCard from "./NotebookCard.svelte";
    import CreateNotebookDialog from "./CreateNotebookDialog.svelte";
    import { toast } from "svelte-sonner";

    let dialogOpen = false;

    $: filteredNotebooks = $notebooksStore.notebooks.filter((notebook) =>
        notebook.name
            .toLowerCase()
            .includes($notebooksStore.searchQuery.toLowerCase()),
    );

    $: paginatedNotebooks = filteredNotebooks.slice(
        ($notebooksStore.currentPage - 1) * $notebooksStore.itemsPerPage,
        $notebooksStore.currentPage * $notebooksStore.itemsPerPage,
    );

    async function handleDelete(id: string) {
        try {
            await notebooksStore.deleteNotebook(id);
            toast.success("Notebook deleted successfully");
        } catch (error) {
            toast.error("Failed to delete notebook");
        }
    }

    onMount(() => {
        notebooksStore.loadNotebooks();
    });
</script>

<div class="space-y-6">
    <div class="flex justify-between items-center">
        <h2 class="text-2xl font-bold">My Notebooks</h2>
        <Button
            on:click={() => (dialogOpen = true)}
            class="flex items-center gap-2"
        >
            <Plus class="w-4 h-4" />
            New Notebook
        </Button>
    </div>

    <div class="flex gap-4 items-center">
        <div class="relative flex-1">
            <div class="absolute left-2 top-2.5 h-4 w-4 text-muted-foreground">
                <Search size={16} />
            </div>
            <Input
                type="text"
                placeholder="Search notebooks..."
                value={$notebooksStore.searchQuery}
                on:input={(e) =>
                    notebooksStore.setSearchQuery(e.currentTarget.value)}
                class="pl-8"
            />
        </div>
    </div>

    {#if $notebooksStore.loading}
        <div class="text-center py-12">
            <p class="text-muted-foreground">Loading notebooks...</p>
        </div>
    {:else}
        <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
            {#each paginatedNotebooks as notebook (notebook.id)}
                <NotebookCard {notebook} onDelete={handleDelete} />
            {/each}
        </div>

        {#if filteredNotebooks.length === 0}
            <div class="text-center py-12">
                <p class="text-muted-foreground">
                    {$notebooksStore.searchQuery
                        ? "No notebooks found matching your search."
                        : "No notebooks yet. Create one to get started!"}
                </p>
            </div>
        {/if}

        {#if $notebooksStore.totalPages > 1}
            <div class="flex justify-center gap-2 mt-4">
                <Button
                    variant="outline"
                    size="icon"
                    disabled={$notebooksStore.currentPage === 1}
                    on:click={() =>
                        notebooksStore.setPage($notebooksStore.currentPage - 1)}
                >
                    <ChevronLeft class="w-4 h-4" />
                </Button>

                <span class="flex items-center px-4 text-sm">
                    Page {$notebooksStore.currentPage} of {$notebooksStore.totalPages}
                </span>

                <Button
                    variant="outline"
                    size="icon"
                    disabled={$notebooksStore.currentPage ===
                        $notebooksStore.totalPages}
                    on:click={() =>
                        notebooksStore.setPage($notebooksStore.currentPage + 1)}
                >
                    <ChevronRight class="w-4 h-4" />
                </Button>
            </div>
        {/if}
    {/if}

    <CreateNotebookDialog bind:open={dialogOpen} />
</div>
