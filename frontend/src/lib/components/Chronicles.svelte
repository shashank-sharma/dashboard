<script lang="ts">
    import { onMount } from "svelte";
    import { Card } from "$lib/components/ui/card";
    import { Button } from "$lib/components/ui/button";
    import { Input } from "$lib/components/ui/input";
    import { Label } from "$lib/components/ui/label";
    import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from "$lib/components/ui/select";
    import { pb } from "$lib/pocketbase";
    import { toast } from "svelte-sonner";
    import { Carta, MarkdownEditor } from "carta-md";
    import ChronicleMetadata from "./ChronicleMetadata.svelte";
    import "carta-md/default.css";
    import DOMPurify from "isomorphic-dompurify";
    import { getContext } from "svelte";

    const { theme } = getContext("theme");

    export let date = new Date();

    $: formattedDate = date.toISOString().split("T")[0].replace(/-/g, "");
    $: displayDate = new Intl.DateTimeFormat("en-US", {
        weekday: "long",
        year: "numeric",
        month: "long",
        day: "numeric",
    }).format(date);

    let content = "";
    let mood = "neutral";
    let tags = "";
    let isLoading = false;
    let editorHeight = 0;

    let carta = new Carta({
        sanitizer: DOMPurify.sanitize,
        plugins: [
            "heading",
            "bold",
            "italic",
            "strikethrough",
            "link",
            "list",
            "table",
            "image",
            "code",
            "blockquote",
        ],
        theme: $theme === "dark" ? "github-dark" : "github-light",
    });

    onMount(async () => {
        try {
            const record = await pb.collection("journal_entries").getFirstListItem(
                `date = "${date.toISOString().split("T")[0]}" && user = "${pb.authStore.model.id}"`
            );

            if (record) {
                content = record.content;
                mood = record.mood;
                tags = record.tags;
            }
        } catch (error) {
            console.log("No existing entry found for today");
        }

        const vh = Math.max(document.documentElement.clientHeight || 0, window.innerHeight || 0);
        editorHeight = vh - 200;
    });

    async function handleSave() {
        if (!content.trim()) {
            toast.error("Please write some content");
            return;
        }

        isLoading = true;

        try {
            const data = {
                user: pb.authStore.model.id,
                title: formattedDate,
                content,
                date: date.toISOString().split("T")[0],
                mood,
                tags,
            };

            try {
                const existingRecord = await pb.collection("journal_entries").getFirstListItem(
                    `date = "${date.toISOString().split("T")[0]}" && user = "${pb.authStore.model.id}"`
                );
                await pb.collection("journal_entries").update(existingRecord.id, data);
                toast.success("Updated successfully");
            } catch {
                await pb.collection("journal_entries").create(data);
                toast.success("Saved successfully");
            }
        } catch (error) {
            console.error("Error saving entry:", error);
            toast.error("Failed to save");
        } finally {
            isLoading = false;
        }
    }
</script>

<div class="container mx-auto p-4">
    <div class="flex gap-4 h-[calc(100vh-100px)]">
        <!-- Main Editor Section -->
        <div class="flex-grow h-full">
            <Card class="h-full flex flex-col">
                <div class="p-4 border-b flex items-center justify-between flex-shrink-0">
                    <div>
                        <h2 class="text-xl font-semibold">{displayDate}</h2>
                        <p class="text-sm text-muted-foreground">Entry #{formattedDate}</p>
                    </div>
                    <div class="flex items-center gap-4">
                        <div class="space-x-2">
                            <Label for="mood" class="text-sm">Mood</Label>
                            <Select bind:value={mood} class="w-32">
                                <SelectTrigger id="mood">
                                    <SelectValue placeholder="Mood" />
                                </SelectTrigger>
                                <SelectContent>
                                    <SelectItem value="happy">Happy</SelectItem>
                                    <SelectItem value="neutral">Neutral</SelectItem>
                                    <SelectItem value="sad">Sad</SelectItem>
                                    <SelectItem value="excited">Excited</SelectItem>
                                    <SelectItem value="anxious">Anxious</SelectItem>
                                    <SelectItem value="peaceful">Peaceful</SelectItem>
                                </SelectContent>
                            </Select>
                        </div>
                        <div class="space-x-2">
                            <Label for="tags" class="text-sm">Tags</Label>
                            <Input 
                                id="tags"
                                bind:value={tags}
                                placeholder="life, work..."
                                class="w-48"
                            />
                        </div>
                        <Button 
                            variant="default"
                            on:click={handleSave}
                            disabled={isLoading}
                        >
                            {isLoading ? "Saving..." : "Save"}
                        </Button>
                    </div>
                </div>

                <!-- Fixed height editor container -->
                <div class="flex-grow overflow-hidden border-b">
                    <div class="h-full" style="height: {editorHeight}px">
                        <MarkdownEditor 
                            {carta} 
                            bind:value={content} 
                            class="carta-editor-custom {$theme === 'dark' ? 'dark' : ''}"
                        />
                    </div>
                </div>
            </Card>
        </div>

        <!-- Metadata Sidebar -->
        <div class="w-96 h-full">
            <ChronicleMetadata />
        </div>
    </div>
</div>

<style>
    :global(.carta-editor) {
        min-height: calc(100vh - 300px);
        max-height: calc(100vh - 300px);
        overflow-y: auto;
        font-family: inherit;
        padding: 1rem !important;
    }

    :global(.carta-editor:focus) {
        outline: none;
    }

    :global(.carta-font-code) {
        font-family: "Menlo", "Monaco", "Courier New", monospace;
        font-size: 1.1rem;
    }

    /* Dark theme overrides */
    :global(.dark .carta-editor),
    :global(.dark .carta-wrapper),
    :global(.dark .carta-container) {
        background-color: hsl(var(--background));
        color: hsl(var(--foreground));
    }

    :global(.dark .carta-toolbar) {
        background-color: hsl(var(--background));
        border-color: hsl(var(--border));
    }

    :global(.dark .carta-toolbar button) {
        color: hsl(var(--foreground));
    }

    :global(.dark .carta-toolbar button:hover) {
        background-color: hsl(var(--accent));
    }

    :global(.dark .carta-input-wrapper) {
        background-color: hsl(var(--background));
    }

    :global(.dark .carta-input-wrapper textarea) {
        background-color: transparent;
        color: hsl(var(--foreground));
    }

    :global(.dark pre),
    :global(.dark code) {
        background-color: hsl(var(--muted));
        color: hsl(var(--muted-foreground));
    }

    :global(.dark blockquote) {
        border-left-color: hsl(var(--border));
        color: hsl(var(--muted-foreground));
    }

    :global(.dark .carta-renderer) {
        background-color: hsl(var(--background));
        color: hsl(var(--foreground));
    }

    /* Override shiki highlighting for dark mode */
    :global(.dark .shiki),
    :global(.dark .shiki span) {
        background-color: hsl(var(--muted)) !important;
        color: hsl(var(--foreground)) !important;
    }

    /* Custom scrollbar */
    :global(.carta-editor-custom ::-webkit-scrollbar) {
        width: 6px;
    }

    :global(.carta-editor-custom ::-webkit-scrollbar-track) {
        background: transparent;
    }

    :global(.carta-editor-custom ::-webkit-scrollbar-thumb) {
        background-color: hsl(var(--muted));
        border-radius: 3px;
    }

    :global(.carta-editor-custom ::-webkit-scrollbar-thumb:hover) {
        background-color: hsl(var(--muted-foreground));
    }
</style>