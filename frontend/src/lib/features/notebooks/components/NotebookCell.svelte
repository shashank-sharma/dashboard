<script lang="ts">
    import { Button } from "$lib/components/ui/button";
    import { Card } from "$lib/components/ui/card";
    import { Trash2, Play, Type } from "lucide-svelte";
    import type { Cell } from "../types";
    import { basicSetup } from "codemirror";
    import { EditorView, keymap } from "@codemirror/view";
    import { EditorState } from "@codemirror/state";
    import { python } from "@codemirror/lang-python";
    import { markdown } from "@codemirror/lang-markdown";
    import { indentUnit } from "@codemirror/language";
    import { HighlightStyle, syntaxHighlighting } from "@codemirror/language";
    import { tags as t } from "@lezer/highlight";
    import { pythonService } from "../services/python.service";
    import { toast } from "svelte-sonner";

    export let cell: Cell;
    export let onUpdate: (updates: Partial<Cell>) => void;
    export let onDelete: () => void;
    export let onToggleType: () => void;

    let isExecuting = false;
    let editorElement: HTMLElement;
    let editor: EditorView;

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
    ]);

    function createEditor(element: HTMLElement) {
        const extensions = [
            basicSetup,
            indentUnit.of("    "),
            EditorView.lineWrapping,
            EditorState.tabSize.of(4),
            cell.type === "code" ? python() : markdown(),
            syntaxHighlighting(customTheme),
            keymap.of([
                {
                    key: "Shift-Enter",
                    run: () => {
                        if (cell.type === "code") {
                            executeCell();
                        }
                        return true;
                    },
                },
            ]),
            EditorView.updateListener.of((update) => {
                if (update.docChanged) {
                    onUpdate({ content: update.state.doc.toString() });
                }
            }),
        ];

        editor = new EditorView({
            doc: cell.content,
            extensions,
            parent: element,
        });
    }

    async function executeCell() {
        if (cell.type !== "code" || isExecuting) return;

        isExecuting = true;
        try {
            const result = await pythonService.executeCode(cell.content);
            onUpdate({ output: result }); // This updates the cell's output
            toast.success("Code executed successfully");
        } catch (error) {
            onUpdate({ output: `Error: ${error.message}` });
            toast.error("Failed to execute code");
        } finally {
            isExecuting = false;
        }
    }
</script>

<Card class="p-4">
    <div class="flex gap-2 mb-2">
        <Button
            size="icon"
            variant="outline"
            on:click={onToggleType}
            title={cell.type === "code" ? "Switch to text" : "Switch to code"}
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
                on:click={executeCell}
                disabled={isExecuting}
                title="Run cell (Shift+Enter)"
            >
                <Play class="w-4 h-4" />
            </Button>
        {/if}

        <Button
            size="icon"
            variant="outline"
            on:click={onDelete}
            title="Delete cell"
        >
            <Trash2 class="w-4 h-4" />
        </Button>
    </div>

    <div
        class="editor-container min-h-[100px] border rounded-md"
        bind:this={editorElement}
        use:createEditor
    />

    {#if cell.type === "code" && cell.output}
        <div class="mt-2 bg-muted p-4 rounded-md">
            <pre
                class="whitespace-pre-wrap font-mono text-sm overflow-x-auto">{cell.output}</pre>
        </div>
    {/if}
</Card>

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

    :global(.editor-container .cm-gutters) {
        background-color: hsl(var(--muted));
        border-right: 1px solid hsl(var(--border));
        color: hsl(var(--muted-foreground));
    }
</style>
