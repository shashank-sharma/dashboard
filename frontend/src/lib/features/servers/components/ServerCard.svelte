<script lang="ts">
    import { Server as ServerIcon, Trash2 } from "lucide-svelte";
    import { Badge } from "$lib/components/ui/badge";
    import { Button } from "$lib/components/ui/button";
    import { Switch } from "$lib/components/ui/switch";
    import {
        Card,
        CardHeader,
        CardTitle,
        CardDescription,
        CardContent,
        CardFooter,
    } from "$lib/components/ui/card";
    import type { Server } from "../types";

    export let server: Server;
    export let onEdit: (server: Server) => void;
    export let onDelete: (id: string) => void;
    export let onToggleStatus: (id: string, status: boolean) => void;

    function handleCardClick(e: MouseEvent) {
        const target = e.target as HTMLElement;
        if (!target.closest(".action-button")) {
            onEdit(server);
        }
    }

    function formatDate(date: string) {
        return new Date(date).toLocaleString();
    }
</script>

<Card
    class="cursor-pointer hover:shadow-md transition-shadow"
    on:click={handleCardClick}
>
    <CardHeader>
        <CardTitle class="flex items-center justify-between">
            <div class="flex items-center gap-2">
                <ServerIcon class="w-4 h-4" />
                <span>{server.name}</span>
            </div>
            <Badge variant={server.is_active ? "default" : "secondary"}>
                {server.is_active ? "Active" : "Disabled"}
            </Badge>
        </CardTitle>
        <CardDescription>{server.url}</CardDescription>
    </CardHeader>
    <CardContent>
        <div class="space-y-2">
            <div class="text-sm">
                <span class="font-medium">Provider:</span>
                {server.provider}
            </div>
            <div class="text-sm">
                <span class="font-medium">Created:</span>
                {formatDate(server.created)}
            </div>
        </div>
    </CardContent>
    <CardFooter class="justify-between">
        <div class="flex items-center gap-2 action-button">
            <Switch
                checked={server.is_active}
                onCheckedChange={() =>
                    onToggleStatus(server.id, server.is_active)}
            />
            <span class="text-sm">
                {server.is_active ? "Active" : "Disabled"}
            </span>
        </div>
        <Button
            variant="destructive"
            size="icon"
            class="action-button"
            on:click={() => onDelete(server.id)}
        >
            <Trash2 class="w-4 h-4" />
        </Button>
    </CardFooter>
</Card>
