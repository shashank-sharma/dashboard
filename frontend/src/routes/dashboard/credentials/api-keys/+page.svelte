<script lang="ts">
    import { onMount } from "svelte";
    import { pb } from "$lib/config/pocketbase";
    import { toast } from "svelte-sonner";
    import { Plus, RefreshCcw } from "lucide-svelte";
    import { Button } from "$lib/components/ui/button";
    import * as AlertDialog from "$lib/components/ui/alert-dialog";
    import type { ApiKey } from "$lib/features/credentials/types";
    import { ApiKeyCard } from "$lib/features/credentials/components";

    let apiKeys: ApiKey[] = [];
    let loading = true;
    let showDeleteDialog = false;
    let keyToDelete: string | null = null;

    async function loadApiKeys() {
        loading = true;
        try {
            const records = await pb.collection("api_keys").getFullList({
                sort: "-created",
                expand: "user",
            });
            apiKeys = records as unknown as ApiKey[];
        } catch (error) {
            toast.error("Failed to load API keys");
        } finally {
            loading = false;
        }
    }

    function handleEdit(apiKey: ApiKey) {
        toast.info("Edit functionality to be implemented");
    }

    function handleDelete(id: string) {
        keyToDelete = id;
        showDeleteDialog = true;
    }

    async function confirmDelete() {
        if (!keyToDelete) return;
        try {
            await pb.collection("api_keys").delete(keyToDelete);
            toast.success("API key deleted successfully");
            loadApiKeys();
        } catch (error) {
            toast.error("Failed to delete API key");
        } finally {
            showDeleteDialog = false;
            keyToDelete = null;
        }
    }

    async function handleToggleStatus(id: string, currentStatus: boolean) {
        try {
            await pb.collection("api_keys").update(id, {
                is_active: !currentStatus,
            });
            toast.success("API key status updated");
            loadApiKeys();
        } catch (error) {
            toast.error("Failed to update API key status");
        }
    }

    onMount(loadApiKeys);
</script>

<div class="p-6">
    <div class="flex justify-between items-center mb-6">
        <h2 class="text-3xl font-bold">API Keys</h2>
        <div class="flex space-x-2">
            <Button variant="outline" on:click={loadApiKeys}>
                <RefreshCcw class="w-4 h-4 mr-2" />
                Refresh
            </Button>
            <Button>
                <Plus class="w-4 h-4 mr-2" />
                New API Key
            </Button>
        </div>
    </div>

    {#if loading}
        <div class="flex justify-center items-center h-64">
            <RefreshCcw class="w-8 h-8 animate-spin" />
        </div>
    {:else if apiKeys.length === 0}
        <div class="flex flex-col items-center justify-center h-64 text-center">
            <p class="text-muted-foreground mb-4">No API keys found</p>
            <Button>
                <Plus class="w-4 h-4 mr-2" />
                Create your first API key
            </Button>
        </div>
    {:else}
        <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
            {#each apiKeys as apiKey (apiKey.id)}
                <ApiKeyCard
                    {apiKey}
                    onEdit={handleEdit}
                    onDelete={handleDelete}
                    onToggleStatus={handleToggleStatus}
                />
            {/each}
        </div>
    {/if}
</div>

<AlertDialog.Root bind:open={showDeleteDialog}>
    <AlertDialog.Content>
        <AlertDialog.Header>
            <AlertDialog.Title>Are you sure?</AlertDialog.Title>
            <AlertDialog.Description>
                This action cannot be undone. This will permanently delete the
                API key.
            </AlertDialog.Description>
        </AlertDialog.Header>
        <AlertDialog.Footer>
            <AlertDialog.Cancel>Cancel</AlertDialog.Cancel>
            <AlertDialog.Action on:click={confirmDelete}
                >Delete</AlertDialog.Action
            >
        </AlertDialog.Footer>
    </AlertDialog.Content>
</AlertDialog.Root>
