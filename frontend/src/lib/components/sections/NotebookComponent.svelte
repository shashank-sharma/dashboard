<script lang="ts">
    import { onMount, onDestroy } from "svelte";
    import { Button } from "$lib/components/ui/button";
    import { Card } from "$lib/components/ui/card";
    import { Alert, AlertDescription } from "$lib/components/ui/alert";
    import { notebookService } from "$lib/services/notebookService";
    import { pyodideService } from "$lib/services/pyodideService";
    import { Play, Plus, Save, Trash2, Type } from "lucide-svelte";
    import { toast } from "svelte-sonner";
    import { basicSetup } from "codemirror";
    import { EditorView } from "@codemirror/view";
    import { HighlightStyle, syntaxHighlighting } from "@codemirror/language";
    import { tags as t } from "@lezer/highlight";

    const customTheme = HighlightStyle.define([
        { tag: t.keyword, color: "var(--code-keyword, #c792ea)" },
        { tag: t.operator, color: "var(--code-operator, #89ddff)" },
        { tag: t.string, color: "var(--code-string, #c3e88d)" },
        { tag: t.number, color: "var(--code-number, #f78c6c)" },
        {
            tag: t.comment,
            color: "var(--code-comment, #546e7a)",
            fontStyle: "italic",
        },
        {
            tag: t.function(t.variableName),
            color: "var(--code-function, #82aaff)",
        },
        {
            tag: t.definition(t.variableName),
            color: "var(--code-variable, #f07178)",
        },
        { tag: t.className, color: "var(--code-class, #ffcb6b)" },
        { tag: t.propertyName, color: "var(--code-property, #80cbc4)" },
        {
            tag: [t.typeName, t.regexp, t.meta],
            color: "var(--code-meta, #c792ea)",
        },
    ]);

    const theme = EditorView.theme({
        "&": {
            backgroundColor: "hsl(var(--background))",
            color: "hsl(var(--foreground))",
        },
        ".cm-content": {
            caretColor: "hsl(var(--foreground))",
        },
        "&.cm-focused .cm-cursor": {
            borderLeftColor: "hsl(var(--foreground))",
        },
        "&.cm-focused .cm-selectionBackground, .cm-selectionBackground": {
            backgroundColor: "hsl(var(--accent) / 0.3)",
        },
        ".cm-activeLine": {
            backgroundColor: "hsl(var(--muted) / 0.3)",
        },
    });
    import { python } from "@codemirror/lang-python";
    import { indentUnit } from "@codemirror/language";
    import { EditorState } from "@codemirror/state";
    import { markdown } from "@codemirror/lang-markdown";

    export let id: string;

    let notebook = null;
    let loading = true;
    let executingCells = new Set();
    let editors = new Map();
    let savingTimeout;

    onMount(async () => {
        try {
            await loadNotebook();
            await pyodideService.initialize();
        } catch (error) {
            toast.error("Failed to initialize notebook");
        }
    });

    onDestroy(() => {
        if (savingTimeout) clearTimeout(savingTimeout);
        editors.forEach((editor) => editor.destroy());
    });

    function createEditor(
        element: HTMLElement,
        content: string,
        type: "code" | "markdown",
    ) {
        const extensions = [
            basicSetup,
            indentUnit.of("    "),
            EditorView.lineWrapping,
            EditorState.tabSize.of(4),
            type === "code" ? python() : markdown(),
            theme,
            syntaxHighlighting(customTheme),
            EditorView.theme({
                "&": { height: "100%" },
                ".cm-scroller": { overflow: "auto" },
                ".cm-content": { minHeight: "100px" },
                "&.cm-focused": { outline: "none" },
            }),
        ];

        const view = new EditorView({
            doc: content,
            extensions,
            parent: element,
            dispatch: (tr) => {
                view.update([tr]);
                if (tr.docChanged) {
                    const cellId = element.dataset.cellId;
                    updateCell(cellId, view.state.doc.toString());
                }
            },
        });
        return view;
    }

    function initializeEditor(element: HTMLElement, cell) {
        if (!element || editors.has(cell.id)) return;

        const editor = createEditor(element, cell.content, cell.type);
        editors.set(cell.id, editor);
    }

    async function loadNotebook() {
        try {
            notebook = await notebookService.getNotebook(id);
        } catch (error) {
            toast.error("Failed to load notebook");
        } finally {
            loading = false;
        }
    }

    async function saveNotebook() {
        try {
            await notebookService.updateNotebook(id, notebook);
            toast.success("Notebook saved");
        } catch (error) {
            toast.error("Failed to save notebook");
        }
    }

    // Debounced save
    function debouncedSave() {
        if (savingTimeout) clearTimeout(savingTimeout);
        savingTimeout = setTimeout(saveNotebook, 2000);
    }

    function addCell(type: "code" | "markdown" = "code") {
        const newCell = {
            id: crypto.randomUUID(),
            content: "",
            output: "",
            type,
            language: type === "code" ? "python" : "markdown",
        };

        notebook.cells = [...notebook.cells, newCell];
        debouncedSave();
    }

    function updateCell(cellId: string, content: string) {
        notebook.cells = notebook.cells.map((cell) =>
            cell.id === cellId ? { ...cell, content } : cell,
        );
        debouncedSave();
    }

    function toggleCellType(cellId: string) {
        notebook.cells = notebook.cells.map((cell) => {
            if (cell.id === cellId) {
                const newType = cell.type === "code" ? "markdown" : "code";
                return {
                    ...cell,
                    type: newType,
                    language: newType === "code" ? "python" : "markdown",
                    output: "",
                };
            }
            return cell;
        });
        debouncedSave();
    }

    async function executeCell(cellId: string) {
        const cell = notebook.cells.find((c) => c.id === cellId);
        if (!cell || cell.type !== "code") return;

        executingCells.add(cellId);
        executingCells = executingCells;

        try {
            const { output, error } = await pyodideService.executeCode(
                cell.content,
            );

            notebook.cells = notebook.cells.map((c) =>
                c.id === cellId
                    ? {
                          ...c,
                          output: error
                              ? `Error: ${error}\n\nOutput: ${output}`
                              : output,
                      }
                    : c,
            );

            if (error) {
                toast.error("Error executing cell");
            }
        } catch (error) {
            toast.error("Failed to execute cell");
            notebook.cells = notebook.cells.map((c) =>
                c.id === cellId ? { ...c, output: error.message } : c,
            );
        } finally {
            executingCells.delete(cellId);
            executingCells = executingCells;
            debouncedSave();
        }
    }

    function deleteCell(cellId: string) {
        notebook.cells = notebook.cells.filter((cell) => cell.id !== cellId);
        debouncedSave();
    }

    function renderMarkdown(content: string) {
        return markdownParser.render(content);
    }

    // Keyboard shortcuts
    function handleKeyDown(event: KeyboardEvent, cellId: string) {
        // Shift + Enter to execute
        if (event.key === "Enter" && event.shiftKey) {
            event.preventDefault();
            const cell = notebook.cells.find((c) => c.id === cellId);
            if (cell.type === "code") {
                executeCell(cellId);
            }
        }
    }
</script>

<div class="max-w-4xl mx-auto space-y-6 p-4">
    {#if loading}
        <div class="text-center py-12">
            <p class="text-muted-foreground">Loading notebook...</p>
        </div>
    {:else}
        <div class="flex justify-between items-center">
            <h1 class="text-2xl font-bold">{notebook.name}</h1>
            <div class="flex gap-2">
                <Button on:click={saveNotebook} class="flex items-center gap-2">
                    <Save class="w-4 h-4" />
                    Save
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
                <Card class="p-4">
                    <div class="flex gap-2 mb-2">
                        <Button
                            size="icon"
                            variant="outline"
                            on:click={() => toggleCellType(cell.id)}
                            title={cell.type === "code"
                                ? "Switch to text"
                                : "Switch to code"}
                        >
                            {#if cell.type === "code"}
                                <Type class="w-4 h-4" />
                            {:else}
                                <Play class="w-4 h-4" />
                            {/if}
                        </Button>

                        {#if cell.type === "code"}
                            <Button
                                size="icon"
                                variant="outline"
                                on:click={() => executeCell(cell.id)}
                                disabled={executingCells.has(cell.id)}
                                title="Run cell"
                            >
                                <Play class="w-4 h-4" />
                            </Button>
                        {/if}

                        <Button
                            size="icon"
                            variant="outline"
                            on:click={() => deleteCell(cell.id)}
                            title="Delete cell"
                        >
                            <Trash2 class="w-4 h-4" />
                        </Button>
                    </div>

                    <div
                        class="editor-container min-h-[100px] border rounded-md"
                        data-cell-id={cell.id}
                        use:initializeEditor={cell}
                    />

                    {#if cell.output}
                        <div class="mt-2 bg-muted p-4 rounded-md">
                            <pre
                                class="whitespace-pre-wrap font-mono text-sm overflow-x-auto">{cell.output}</pre>
                        </div>
                    {/if}
                </Card>
            {/each}
        </div>

        {#if notebook.cells.length === 0}
            <Alert>
                <AlertDescription>
                    No cells in this notebook. Add a cell to get started!
                </AlertDescription>
            </Alert>
        {/if}
    {/if}
</div>

<style>
    :global(.editor-container .cm-editor) {
        height: 100%;
        min-height: 100px;
        padding: 8px;
    }

    :global(.editor-container .cm-editor.cm-focused) {
        outline: none;
        border-color: hsl(var(--ring));
    }

    :global(.editor-container .cm-line) {
        padding: 0;
    }

    :global(.editor-container .cm-gutters) {
        background-color: hsl(var(--muted));
        border-right: 1px solid hsl(var(--border));
        color: hsl(var(--muted-foreground));
        margin-right: 10px;
    }

    :global(.editor-container .cm-activeLineGutter) {
        background-color: hsl(var(--accent));
        color: hsl(var(--accent-foreground));
    }
</style>
