<script lang="ts">
    import { fade } from "svelte/transition";
    import { Edit, Trash2, Eye, EyeOff, Power } from "lucide-svelte";
    import { Button } from "$lib/components/ui/button";
    import * as Card from "$lib/components/ui/card";
    import { Badge } from "$lib/components/ui/badge";
    import type { ApiKey } from "../types";
    import { formatDistanceToNow } from "date-fns";

    export let apiKey: ApiKey;
    export let onEdit: (apiKey: ApiKey) => void;
    export let onDelete: (id: string) => void;
    export let onToggleStatus: (id: string, isActive: boolean) => void;

    let showKey = false;
    let showSecret = false;

    function formatDate(dateString: string) {
        return formatDistanceToNow(new Date(dateString), { addSuffix: true });
    }
</script>

<div transition:fade={{ duration: 200 }}>
    <Card.Root class="h-full">
        <Card.Header>
            <div class="flex items-center justify-between">
                <Card.Title>{apiKey.service}</Card.Title>
                <Badge
                    variant={apiKey.is_active ? "default" : "destructive"}
                    class="ml-2"
                >
                    {apiKey.is_active ? "Active" : "Inactive"}
                </Badge>
            </div>
        </Card.Header>
        <Card.Content>
            <div class="space-y-2">
                <div>
                    <p class="text-sm font-medium">API Key</p>
                    <div class="flex items-center mt-1">
                        <p class="text-sm font-mono truncate flex-1">
                            {showKey ? apiKey.key : "••••••••••••••••"}
                        </p>
                        <Button
                            variant="ghost"
                            size="sm"
                            on:click={() => (showKey = !showKey)}
                        >
                            <svelte:component
                                this={showKey ? EyeOff : Eye}
                                class="h-4 w-4"
                            />
                        </Button>
                    </div>
                </div>
                {#if apiKey.secret}
                    <div>
                        <p class="text-sm font-medium">Secret</p>
                        <div class="flex items-center mt-1">
                            <p class="text-sm font-mono truncate flex-1">
                                {showSecret
                                    ? apiKey.secret
                                    : "••••••••••••••••"}
                            </p>
                            <Button
                                variant="ghost"
                                size="sm"
                                on:click={() => (showSecret = !showSecret)}
                            >
                                <svelte:component
                                    this={showSecret ? EyeOff : Eye}
                                    class="h-4 w-4"
                                />
                            </Button>
                        </div>
                    </div>
                {/if}
                <div>
                    <p class="text-sm font-medium">Created</p>
                    <p class="text-sm text-muted-foreground">
                        {formatDate(apiKey.created)}
                    </p>
                </div>
                {#if apiKey.expires}
                    <div>
                        <p class="text-sm font-medium">Expires</p>
                        <p class="text-sm text-muted-foreground">
                            {formatDate(apiKey.expires)}
                        </p>
                    </div>
                {/if}
            </div>
        </Card.Content>
        <Card.Footer class="flex justify-between">
            <Button
                variant="ghost"
                size="sm"
                on:click={() => onToggleStatus(apiKey.id, apiKey.is_active)}
            >
                <Power class="h-4 w-4 mr-2" />
                {apiKey.is_active ? "Deactivate" : "Activate"}
            </Button>
            <div class="space-x-2">
                <Button
                    variant="ghost"
                    size="sm"
                    on:click={() => onEdit(apiKey)}
                >
                    <Edit class="h-4 w-4" />
                </Button>
                <Button
                    variant="ghost"
                    size="sm"
                    on:click={() => onDelete(apiKey.id)}
                >
                    <Trash2 class="h-4 w-4" />
                </Button>
            </div>
        </Card.Footer>
    </Card.Root>
</div>
