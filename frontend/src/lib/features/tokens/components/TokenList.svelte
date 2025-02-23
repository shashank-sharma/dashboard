<script lang="ts">
    import { onMount } from "svelte";
    import { pb } from "$lib/config/pocketbase";
    import { toast } from "svelte-sonner";
    import { Plus, RefreshCcw } from "lucide-svelte";
    import { Button } from "$lib/components/ui/button";
    import TokenCard from "./TokenCard.svelte";
    import TokenDialog from "./TokenDialog.svelte";
    import * as AlertDialog from "$lib/components/ui/alert-dialog";
    import type { Token } from "../types";

    let tokens: Token[] = [];
    let loading = true;
    let showDialog = false;
    let showDeleteDialog = false;
    let selectedToken: Token | null = null;
    let tokenToDelete: string | null = null;

    async function loadTokens() {
        loading = true;
        try {
            const records = await pb.collection("tokens").getFullList({
                sort: "-created",
                expand: "user",
            });
            tokens = records;
        } catch (error) {
            console.error("Load tokens error:", error);
            toast.error("Failed to load tokens");
        } finally {
            loading = false;
        }
    }

    async function handleTokenSubmit(data: any) {
        try {
            console.log("Submitting token data:", data); // Debug log
            if (selectedToken) {
                await pb.collection("tokens").update(selectedToken.id, data);
                toast.success("Token updated successfully");
            } else {
                const tokenData = {
                    ...data,
                    user: pb.authStore.model?.id,
                    is_active: true,
                };
                await pb.collection("tokens").create(tokenData);
                toast.success("Token created successfully");
            }
            showDialog = false;
            await loadTokens();
        } catch (error) {
            console.error("Token submission error:", error);
            toast.error(
                selectedToken
                    ? "Failed to update token"
                    : "Failed to create token",
            );
        }
    }

    function handleEdit(token: Token) {
        selectedToken = token;
        showDialog = true;
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

<div class="container mx-auto p-4">
    <div class="flex justify-between items-center mb-6">
        <h2 class="text-3xl font-bold">API Tokens</h2>
        <Button
            on:click={() => {
                selectedToken = null;
                showDialog = true;
            }}
        >
            <Plus class="w-4 h-4 mr-2" />
            New Token
        </Button>
    </div>

    {#if loading}
        <div class="flex justify-center items-center h-64">
            <RefreshCcw class="w-8 h-8 animate-spin" />
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

    <TokenDialog
        bind:open={showDialog}
        {selectedToken}
        onClose={() => {
            showDialog = false;
            selectedToken = null;
        }}
        onSubmit={handleTokenSubmit}
    />

    <AlertDialog.Root bind:open={showDeleteDialog}>
        <AlertDialog.Content>
            <AlertDialog.Header>
                <AlertDialog.Title>Are you sure?</AlertDialog.Title>
                <AlertDialog.Description>
                    This action cannot be undone. This will permanently delete
                    the token.
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
</div>
