<script lang="ts">
    import { fade } from "svelte/transition";
    import { Edit, Trash2, Eye, EyeOff, Power } from "lucide-svelte";
    import { Button } from "$lib/components/ui/button";
    import * as Card from "$lib/components/ui/card";
    import { Badge } from "$lib/components/ui/badge";
    import type { DeveloperToken } from "../types";
    import { formatDistanceToNow } from "date-fns";

    export let token: DeveloperToken;
    export let onEdit: (token: DeveloperToken) => void;
    export let onDelete: (id: string) => void;
    export let onToggleStatus: (id: string, isActive: boolean) => void;

    let showToken = false;

    function formatDate(dateString: string) {
        return formatDistanceToNow(new Date(dateString), { addSuffix: true });
    }
</script>

<div transition:fade={{ duration: 200 }}>
    <Card.Root class="h-full">
        <Card.Header>
            <div class="flex items-center justify-between">
                <Card.Title>{token.name}</Card.Title>
                <Badge
                    variant={token.is_active ? "default" : "destructive"}
                    class="ml-2"
                >
                    {token.is_active ? "Active" : "Inactive"}
                </Badge>
            </div>
            <Card.Description>
                Environment: {token.environment}
            </Card.Description>
        </Card.Header>
        <Card.Content>
            <div class="space-y-2">
                <div>
                    <p class="text-sm font-medium">Token</p>
                    <div class="flex items-center mt-1">
                        <p class="text-sm font-mono truncate flex-1">
                            {showToken ? token.token : "••••••••••••••••"}
                        </p>
                        <Button
                            variant="ghost"
                            size="sm"
                            on:click={() => (showToken = !showToken)}
                        >
                            <svelte:component
                                this={showToken ? EyeOff : Eye}
                                class="h-4 w-4"
                            />
                        </Button>
                    </div>
                </div>
                <div>
                    <p class="text-sm font-medium">Created</p>
                    <p class="text-sm text-muted-foreground">
                        {formatDate(token.created)}
                    </p>
                </div>
                {#if token.expires}
                    <div>
                        <p class="text-sm font-medium">Expires</p>
                        <p class="text-sm text-muted-foreground">
                            {formatDate(token.expires)}
                        </p>
                    </div>
                {/if}
            </div>
        </Card.Content>
        <Card.Footer class="flex justify-between">
            <Button
                variant="ghost"
                size="sm"
                on:click={() => onToggleStatus(token.id, token.is_active)}
            >
                <Power class="h-4 w-4 mr-2" />
                {token.is_active ? "Deactivate" : "Activate"}
            </Button>
            <div class="space-x-2">
                <Button
                    variant="ghost"
                    size="sm"
                    on:click={() => onEdit(token)}
                >
                    <Edit class="h-4 w-4" />
                </Button>
                <Button
                    variant="ghost"
                    size="sm"
                    on:click={() => onDelete(token.id)}
                >
                    <Trash2 class="h-4 w-4" />
                </Button>
            </div>
        </Card.Footer>
    </Card.Root>
</div>
