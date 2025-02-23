<script lang="ts">
    import { onMount } from "svelte";
    import { notebooksService } from "../services/notebooks.service";
    import { pythonService } from "../services/python.service";
    import { Button } from "$lib/components/ui/button";
    import { Input } from "$lib/components/ui/input";
    import { Save, Plus, Type } from "lucide-svelte";
    import { toast } from "svelte-sonner";
    import type { Notebook, Cell } from "../types";
    import NotebookCell from "./NotebookCell.svelte";

    export let id: string;

    let notebook: Notebook | null = null;
    let loading = true;
    let saving = false;

    async function loadNotebook() {
        try {
            notebook = await notebooksService.getNotebook(id);
        } catch (error) {
            toast.error("Failed to load notebook");
        } finally {
            loading = false;
        }
    }

    async function saveNotebook() {
        if (!notebook) return;

        saving = true;
        try {
            await notebooksService.updateNotebook(id, notebook);
            toast.success("Notebook saved successfully");
        } catch (error) {
            toast.error("Failed to save notebook");
        } finally {
            saving = false;
        }
    }

    function addCell(type: "code" | "markdown") {
        if (!notebook) return;

        const newCell: Cell = {
            id: crypto.randomUUID(),
            type,
            content: "",
            output: "",
            language: type === "code" ? "python" : "markdown",
        };

        notebook.cells = [...notebook.cells, newCell];
    }

    async function executeCell(cellId: string) {
        const cell = notebook?.cells.find((c) => c.id === cellId);
        if (!cell || cell.type !== "code") return;

        try {
            const output = await pythonService.executeCode(cell.content);
            updateCell(cellId, { output });
        } catch (error) {
            updateCell(cellId, { output: `Error: ${error.message}` });
            throw error;
        }
    }

    function updateCell(cellId: string, updates: Partial<Cell>) {
        if (!notebook) return;
        notebook.cells = notebook.cells.map((cell) =>
            cell.id === cellId ? { ...cell, ...updates } : cell,
        );
        debouncedSave();
    }

    let savingTimeout: NodeJS.Timeout;
    function debouncedSave() {
        if (savingTimeout) clearTimeout(savingTimeout);
        savingTimeout = setTimeout(() => saveNotebook(), 2000);
    }

    onMount(async () => {
        try {
            await loadNotebook();
            await pythonService.initialize();
        } catch (error) {
            toast.error("Failed to initialize notebook");
        } finally {
            loading = false;
        }
    });

    function deleteCell(cellId: string) {
        if (!notebook) return;

        notebook.cells = notebook.cells.filter((cell) => cell.id !== cellId);
    }

    function toggleCellType(cellId: string) {
        if (!notebook) return;

        notebook.cells = notebook.cells.map((cell) =>
            cell.id === cellId
                ? {
                      ...cell,
                      type: cell.type === "code" ? "markdown" : "code",
                      language: cell.type === "code" ? "markdown" : "python",
                      output: "",
                  }
                : cell,
        );
    }
</script>

<div class="max-w-4xl mx-auto space-y-6 p-4">
    {#if loading}
        <div class="text-center py-12">
            <p class="text-muted-foreground">Loading notebook...</p>
        </div>
    {:else if notebook}
        <div class="flex justify-between items-center">
            <Input
                type="text"
                bind:value={notebook.name}
                class="text-2xl font-bold bg-transparent border-none h-auto p-0 focus-visible:ring-0"
            />

            <div class="flex gap-2">
                <Button
                    on:click={saveNotebook}
                    disabled={saving}
                    class="flex items-center gap-2"
                >
                    <Save class="w-4 h-4" />
                    {saving ? "Saving..." : "Save"}
                </Button>

                <div class="flex gap-2">
                    <Button
                        on:click={() => addCell("code")}
                        variant="outline"
                        class="flex items-center gap-2"
                    >
                        <Plus class="w-4 h-4" />
                        Code
                    </Button>
                    <Button
                        on:click={() => addCell("markdown")}
                        variant="outline"
                        class="flex items-center gap-2"
                    >
                        <Type class="w-4 h-4" />
                        Text
                    </Button>
                </div>
            </div>
        </div>

        <div class="space-y-4">
            {#each notebook.cells as cell (cell.id)}
                <NotebookCell
                    {cell}
                    onUpdate={(updates) => updateCell(cell.id, updates)}
                    onDelete={() => deleteCell(cell.id)}
                    onToggleType={() => toggleCellType(cell.id)}
                    onExecute={executeCell}
                />
            {/each}
        </div>
    {/if}
</div>
