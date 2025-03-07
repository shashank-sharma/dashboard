<script lang="ts">
    import { onMount } from "svelte";
    import { pb } from "$lib/config/pocketbase";
    import { toast } from "svelte-sonner";
    import {
        Plus,
        RefreshCcw,
        KeyRound,
        Eye,
        EyeOff,
        Copy,
        Check,
        Pencil,
        Trash2,
    } from "lucide-svelte";
    import { Button } from "$lib/components/ui/button";
    import {
        Card,
        CardContent,
        CardDescription,
        CardFooter,
        CardHeader,
        CardTitle,
    } from "$lib/components/ui/card";
    import { Badge } from "$lib/components/ui/badge";
    import { credentialsStore } from "$lib/features/credentials/stores/credentials.store";
    // @ts-ignore
    import SecurityKeyDialog from "./SecurityKeyDialog.svelte";
    import * as AlertDialog from "$lib/components/ui/alert-dialog";
    import type { SecurityKey } from "$lib/features/credentials/types";
    import { Switch } from "$lib/components/ui/switch";

    let securityKeys: SecurityKey[] = [];
    let loading = true;
    let showDialog = false;
    let showDeleteDialog = false;
    let selectedKey: SecurityKey | null = null;
    let keyToDelete: string | null = null;
    let showingPrivateKey: Record<string, boolean> = {};

    async function loadSecurityKeys() {
        loading = true;
        try {
            const records = await pb.collection("security_keys").getFullList({
                sort: "-created",
                expand: "user",
            });
            securityKeys = records as unknown as SecurityKey[];

            // Update only the security keys count without fetching other data
            credentialsStore.updateSecurityKeysCount(securityKeys.length);
        } catch (error: any) {
            console.error("Load security keys error:", error);
            toast.error("Failed to load security keys");
        } finally {
            loading = false;
        }
    }

    async function handleKeySubmit(data: any) {
        try {
            // Create a clean submission object with proper type handling
            const submissionData = {
                name: data.name,
                description: data.description || "",
                private_key: data.private_key,
                public_key: data.public_key,
                is_active:
                    data.is_active === undefined ? true : !!data.is_active,
            };

            if (selectedKey) {
                await pb
                    .collection("security_keys")
                    .update(selectedKey.id, submissionData);
                toast.success("Security key updated successfully");
            } else {
                const keyData = {
                    ...submissionData,
                    user: pb.authStore.model?.id,
                };
                await pb.collection("security_keys").create(keyData);
                toast.success("Security key created successfully");
            }
            showDialog = false;
            selectedKey = null; // Clear selected key after submission
            await loadSecurityKeys();
        } catch (error: any) {
            console.error("Security key submission error:", error);
            toast.error(
                selectedKey
                    ? "Failed to update security key"
                    : "Failed to create security key",
            );
        }
    }

    function handleEdit(key: SecurityKey) {
        // Ensure proper type handling for is_active
        const is_active = key.is_active === undefined ? true : !!key.is_active;

        // Create a complete copy with all needed fields
        selectedKey = {
            id: key.id,
            name: key.name || "",
            description: key.description || "",
            private_key: key.private_key || "",
            public_key: key.public_key || "",
            is_active: is_active,
            created: key.created,
            updated: key.updated,
            user: key.user,
        };

        showDialog = true;
    }

    function handleDelete(id: string) {
        keyToDelete = id;
        showDeleteDialog = true;
    }

    async function confirmDelete() {
        if (!keyToDelete) return;
        try {
            await pb.collection("security_keys").delete(keyToDelete);
            toast.success("Security key deleted successfully");
            await loadSecurityKeys();
        } catch (error) {
            toast.error("Failed to delete security key");
        } finally {
            showDeleteDialog = false;
            keyToDelete = null;
        }
    }

    function formatDate(date: string) {
        return new Date(date).toLocaleString();
    }

    function togglePrivateKeyVisibility(id: string) {
        showingPrivateKey = {
            ...showingPrivateKey,
            [id]: !showingPrivateKey[id],
        };
    }

    async function copyToClipboard(text: string, type: string) {
        try {
            await navigator.clipboard.writeText(text);
            toast.success(`${type} copied to clipboard`);
        } catch (error) {
            toast.error(`Could not copy ${type}`);
        }
    }

    onMount(loadSecurityKeys);
</script>

<div class="container mx-auto px-2 sm:px-4">
    <div
        class="flex flex-col sm:flex-row justify-between items-start sm:items-center mb-6 gap-4"
    >
        <h2 class="text-2xl sm:text-3xl font-bold">Security Keys</h2>
        <div class="flex gap-2 w-full sm:w-auto">
            <Button
                variant="outline"
                class="flex-1 sm:flex-initial"
                on:click={loadSecurityKeys}
            >
                <RefreshCcw class="w-4 h-4 mr-2" />
                Refresh
            </Button>
            <Button
                class="flex-1 sm:flex-initial"
                on:click={() => {
                    selectedKey = null;
                    showDialog = true;
                }}
            >
                <Plus class="w-4 h-4 mr-2" />
                New Key
            </Button>
        </div>
    </div>

    {#if loading}
        <div class="flex justify-center items-center h-64">
            <RefreshCcw class="w-8 h-8 animate-spin" />
        </div>
    {:else if securityKeys.length === 0}
        <div class="text-center py-8">
            <KeyRound class="w-12 h-12 mx-auto text-muted-foreground" />
            <p class="mt-4 text-muted-foreground">No security keys found</p>
            <Button
                variant="outline"
                class="mt-4"
                on:click={() => {
                    selectedKey = null;
                    showDialog = true;
                }}
            >
                <Plus class="w-4 h-4 mr-2" />
                Add Your First Key
            </Button>
        </div>
    {:else}
        <div
            class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-3"
        >
            {#each securityKeys as key (key.id)}
                <Card
                    class="overflow-hidden border-opacity-50 hover:border-opacity-100 transition-all h-full"
                >
                    <CardHeader class="px-3 py-2">
                        <div class="flex flex-col gap-1">
                            <div class="flex items-center justify-between">
                                <div class="overflow-hidden">
                                    <CardTitle
                                        class="text-sm flex items-center gap-1 truncate"
                                    >
                                        <KeyRound
                                            class="w-3.5 h-3.5 flex-shrink-0"
                                        />
                                        <span class="truncate">{key.name}</span>
                                    </CardTitle>
                                </div>
                                {#if key.is_active}
                                    <Badge
                                        variant="outline"
                                        class="px-1.5 py-0 border-green-200 bg-green-50 text-green-700 text-[10px] font-normal flex items-center gap-1 shrink-0"
                                    >
                                        <div
                                            class="w-1.5 h-1.5 rounded-full bg-green-500"
                                        ></div>
                                        Active
                                    </Badge>
                                {:else}
                                    <Badge
                                        variant="outline"
                                        class="px-1.5 py-0 border-gray-200 bg-gray-50 text-gray-500 text-[10px] font-normal flex items-center gap-1 shrink-0"
                                    >
                                        <div
                                            class="w-1.5 h-1.5 rounded-full bg-gray-400"
                                        ></div>
                                        Inactive
                                    </Badge>
                                {/if}
                            </div>
                            <CardDescription class="text-xs line-clamp-1">
                                {key.description || "No description"}
                            </CardDescription>
                        </div>
                    </CardHeader>
                    <CardContent class="px-3 py-1">
                        <div class="space-y-1">
                            <div>
                                <h4
                                    class="text-xs font-medium mb-1 flex items-center gap-1"
                                >
                                    <span>Public Key</span>
                                </h4>
                                <div class="relative">
                                    <pre
                                        class="bg-muted p-1.5 rounded-md text-[10px] overflow-x-auto max-h-[40px]">{key.public_key}</pre>
                                    <Button
                                        size="icon"
                                        variant="ghost"
                                        class="absolute top-0.5 right-0.5 h-5 w-5 action-button"
                                        title="Copy public key"
                                        on:click={() =>
                                            copyToClipboard(
                                                key.public_key,
                                                "Public key",
                                            )}
                                    >
                                        <Copy class="h-2.5 w-2.5" />
                                    </Button>
                                </div>
                            </div>
                            <div>
                                <div
                                    class="flex justify-between items-center mb-1"
                                >
                                    <h4
                                        class="text-xs font-medium flex items-center gap-1"
                                    >
                                        <span>Private Key</span>
                                    </h4>
                                    <Button
                                        variant="ghost"
                                        size="sm"
                                        class="h-5 px-1 action-button"
                                        on:click={() =>
                                            togglePrivateKeyVisibility(key.id)}
                                    >
                                        {#if showingPrivateKey[key.id]}
                                            <EyeOff class="w-2.5 h-2.5" />
                                        {:else}
                                            <Eye class="w-2.5 h-2.5" />
                                        {/if}
                                    </Button>
                                </div>
                                {#if showingPrivateKey[key.id]}
                                    <div class="relative">
                                        <pre
                                            class="bg-muted p-1.5 rounded-md text-[10px] overflow-x-auto max-h-[70px]">{key.private_key}</pre>
                                        <Button
                                            size="icon"
                                            variant="ghost"
                                            class="absolute top-0.5 right-0.5 h-5 w-5 action-button"
                                            title="Copy private key"
                                            on:click={() =>
                                                copyToClipboard(
                                                    key.private_key,
                                                    "Private key",
                                                )}
                                        >
                                            <Copy class="h-2.5 w-2.5" />
                                        </Button>
                                    </div>
                                {:else}
                                    <div
                                        class="bg-muted p-1.5 rounded-md text-[10px] flex items-center justify-center h-6"
                                    >
                                        <span
                                            class="text-muted-foreground text-[10px]"
                                            >Hidden for security</span
                                        >
                                    </div>
                                {/if}
                            </div>
                            <div class="text-[10px] text-muted-foreground mt-2">
                                Created: {formatDate(key.created)}
                            </div>
                        </div>
                    </CardContent>
                    <CardFooter
                        class="p-2 flex gap-1 border-t border-border/50 justify-end"
                    >
                        <Button
                            variant="ghost"
                            size="icon"
                            class="h-7 w-7 action-button"
                            title="Edit"
                            on:click={() => handleEdit(key)}
                        >
                            <Pencil class="w-3 h-3" />
                        </Button>
                        <Button
                            variant="ghost"
                            size="icon"
                            class="h-7 w-7 action-button text-destructive hover:text-destructive"
                            title="Delete"
                            on:click={() => handleDelete(key.id)}
                        >
                            <Trash2 class="w-3 h-3" />
                        </Button>
                    </CardFooter>
                </Card>
            {/each}
        </div>
    {/if}

    <SecurityKeyDialog
        bind:open={showDialog}
        {selectedKey}
        onClose={() => {
            showDialog = false;
            selectedKey = null;
        }}
        onSubmit={handleKeySubmit}
    />

    <AlertDialog.Root bind:open={showDeleteDialog}>
        <AlertDialog.Content class="max-w-[95vw] w-full sm:max-w-md">
            <AlertDialog.Header>
                <AlertDialog.Title>Are you sure?</AlertDialog.Title>
                <AlertDialog.Description>
                    This action cannot be undone. This will permanently delete
                    your security key and remove it from our servers.
                </AlertDialog.Description>
            </AlertDialog.Header>
            <AlertDialog.Footer class="flex-col sm:flex-row gap-2">
                <AlertDialog.Cancel class="w-full sm:w-auto"
                    >Cancel</AlertDialog.Cancel
                >
                <AlertDialog.Action
                    on:click={confirmDelete}
                    class="w-full sm:w-auto">Delete</AlertDialog.Action
                >
            </AlertDialog.Footer>
        </AlertDialog.Content>
    </AlertDialog.Root>
</div>

<style>
    /* Enhanced mobile styles */
    @media (max-width: 640px) {
        :global(.action-button) {
            transform: scale(1.1);
        }

        :global(pre) {
            font-size: 9px;
        }
    }
</style>
