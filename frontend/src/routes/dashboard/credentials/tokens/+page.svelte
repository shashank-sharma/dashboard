<script lang="ts">
    import { onMount } from "svelte";
    import { pb } from "$lib/config/pocketbase";
    import { toast } from "svelte-sonner";
    import { Plus, RefreshCcw } from "lucide-svelte";
    import { Button } from "$lib/components/ui/button";
    import * as AlertDialog from "$lib/components/ui/alert-dialog";
    import type { Token } from "$lib/features/tokens/types";
    import { TokenCard } from "$lib/features/tokens/components";

    let tokens: Token[] = [];
    let loading = true;
    let showDeleteDialog = false;
    let tokenToDelete: string | null = null;

    async function loadTokens() {
        loading = true;
        try {
            const records = await pb.collection("tokens").getFullList({
                sort: "-created",
                expand: "user",
            });
            tokens = records as unknown as Token[];
        } catch (error) {
            console.error("Load tokens error:", error);
            toast.error("Failed to load tokens");
        } finally {
            loading = false;
        }
    }

    function handleEdit(token: Token) {
        // Implement edit functionality
        toast.info("Edit functionality to be implemented");
    }

    function handleDelete(id: string) {
        tokenToDelete = id;
        showDeleteDialog = true;
    }

    async function confirmDelete() {
        if (!tokenToDelete) return;
        try {
            await pb.collection("tokens").delete(tokenToDelete);
            toast.success("Token deleted successfully");
            loadTokens();
        } catch (error) {
            toast.error("Failed to delete token");
        } finally {
            showDeleteDialog = false;
            tokenToDelete = null;
        }
    }

    async function handleToggleStatus(id: string, currentStatus: boolean) {
        try {
            await pb.collection("tokens").update(id, {
                is_active: !currentStatus,
            });
            toast.success("Token status updated");
            loadTokens();
        } catch (error) {
            toast.error("Failed to update token status");
        }
    }

    onMount(loadTokens);
</script>

<div class="p-6">
    <div class="flex justify-between items-center mb-6">
        <h2 class="text-3xl font-bold">API Tokens</h2>
        <div class="flex space-x-2">
            <Button variant="outline" on:click={loadTokens}>
                <RefreshCcw class="w-4 h-4 mr-2" />
                Refresh
            </Button>
            <Button>
                <Plus class="w-4 h-4 mr-2" />
                New Token
            </Button>
        </div>
    </div>

    {#if loading}
        <div class="flex justify-center items-center h-64">
            <RefreshCcw class="w-8 h-8 animate-spin" />
        </div>
    {:else if tokens.length === 0}
        <div class="flex flex-col items-center justify-center h-64 text-center">
            <p class="text-muted-foreground mb-4">No API tokens found</p>
            <Button>
                <Plus class="w-4 h-4 mr-2" />
                Create your first token
            </Button>
        </div>
    {:else}
        <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
            {#each tokens as token (token.id)}
                <TokenCard
                    {token}
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
                token.
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
