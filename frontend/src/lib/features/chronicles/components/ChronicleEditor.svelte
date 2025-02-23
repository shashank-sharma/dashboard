<script lang="ts">
    import { chroniclesStore } from "../stores";
    import { Card } from "$lib/components/ui/card";
    import { Button } from "$lib/components/ui/button";
    import { Input } from "$lib/components/ui/input";
    import { Label } from "$lib/components/ui/label";
    import {
        Select,
        SelectContent,
        SelectItem,
        SelectTrigger,
        SelectValue,
    } from "$lib/components/ui/select";
    import { Carta, MarkdownEditor } from "carta-md";
    import "carta-md/default.css";
    import DOMPurify from "isomorphic-dompurify";
    import { getContext } from "svelte";

    const { theme } = getContext("theme");

    let content = $chroniclesStore.currentEntry?.content || "";
    let mood = $chroniclesStore.currentEntry?.mood || "neutral";
    let tags = $chroniclesStore.currentEntry?.tags || "";

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

    async function handleSave() {
        if (!content.trim()) {
            return;
        }

        await chroniclesStore.saveEntry({
            content,
            mood,
            tags,
            date: new Date().toISOString().split("T")[0],
            title: new Date().toISOString().split("T")[0].replace(/-/g, ""),
            user: "", // Will be set in store
        });
    }
</script>

<!-- Editor UI implementation -->
