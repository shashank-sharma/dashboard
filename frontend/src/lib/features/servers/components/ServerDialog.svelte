<script lang="ts">
    import { onMount } from "svelte";
    import { KeyRound, ServerIcon, Terminal } from "lucide-svelte";
    import * as Dialog from "$lib/components/ui/dialog";
    import { Button } from "$lib/components/ui/button";
    import { Input } from "$lib/components/ui/input";
    import { Label } from "$lib/components/ui/label";
    import { Switch } from "$lib/components/ui/switch";
    import * as Select from "$lib/components/ui/select";
    import { PROVIDERS } from "../constants";
    import { DEFAULT_SERVER_FORM } from "../types";
    import type { Server } from "../types";
    import { toast } from "svelte-sonner";
    import { pb } from "$lib/config/pocketbase";
    import type { SecurityKey } from "$lib/features/credentials/types";

    export let open = false;
    export let onClose: () => void;
    export let onSubmit: (data: any) => Promise<void>;
    export let selectedServer: Server | null = null;

    let formData = { ...DEFAULT_SERVER_FORM };
    let isSubmitting = false;
    let securityKeys: SecurityKey[] = [];
    let loadingKeys = false;

    $: if (open) {
        if (selectedServer) {
            formData = { ...selectedServer };
            loadSecurityKeys();
        } else {
            formData = { ...DEFAULT_SERVER_FORM };
            loadSecurityKeys();
        }
    }

    // Load all available security keys
    async function loadSecurityKeys() {
        loadingKeys = true;
        try {
            const records = await pb.collection("security_keys").getFullList({
                sort: "-created",
                filter: "is_active=true",
            });
            securityKeys = records as unknown as SecurityKey[];
        } catch (error) {
            toast.error("Failed to load security keys");
        } finally {
            loadingKeys = false;
        }
    }

    async function handleSubmit() {
        if (!formData.name || !formData.provider || !formData.ip) {
            toast.error("Please fill in all required fields");
            return;
        }

        if (formData.ssh_enabled) {
            if (!formData.username || !formData.security_key) {
                toast.error("Please fill in all required SSH fields");
                return;
            }
        }

        isSubmitting = true;
        try {
            await onSubmit({
                ...formData,
                is_reachable: formData.ssh_enabled,
            });
        } catch (error) {
            toast.error("Failed to submit server");
        } finally {
            isSubmitting = false;
        }
    }

    onMount(() => {
        if (open) {
            loadSecurityKeys();
        }
    });
</script>

<Dialog.Root bind:open onOpenChange={onClose}>
    <Dialog.Content class="sm:max-w-[500px]">
        <Dialog.Header>
            <Dialog.Title>
                {selectedServer ? "Edit Server" : "Add New Server"}
            </Dialog.Title>
            <Dialog.Description>
                {selectedServer
                    ? "Modify your server details"
                    : "Add a new server for SSH connection"}
            </Dialog.Description>
        </Dialog.Header>

        <form class="space-y-4 mt-4" on:submit|preventDefault={handleSubmit}>
            <div class="space-y-4">
                <div class="grid grid-cols-1 gap-4">
                    <div class="space-y-2">
                        <Label for="name" class="flex items-center gap-1">
                            <ServerIcon class="h-3.5 w-3.5" />
                            <span>Server Name *</span>
                        </Label>
                        <Input
                            id="name"
                            bind:value={formData.name}
                            placeholder="Enter server name"
                            required
                        />
                    </div>

                    <div class="space-y-2">
                        <Label for="provider">Provider *</Label>
                        <div class="relative w-full">
                            <select
                                id="provider"
                                class="w-full h-10 rounded-md border border-input bg-background px-3 py-2 text-sm ring-offset-background focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2 disabled:cursor-not-allowed disabled:opacity-50"
                                bind:value={formData.provider}
                            >
                                <option value="" disabled
                                    >Select provider</option
                                >
                                {#each PROVIDERS as provider}
                                    <option value={provider.value}
                                        >{provider.label}</option
                                    >
                                {/each}
                            </select>
                        </div>
                    </div>

                    <div class="space-y-2">
                        <Label for="ip">IP Address *</Label>
                        <Input
                            id="ip"
                            bind:value={formData.ip}
                            placeholder="Enter server IP address"
                            required
                        />
                    </div>
                </div>
            </div>

            <div class="pt-2">
                <div class="flex items-center justify-between mb-4">
                    <div class="flex items-center gap-2">
                        <Terminal class="h-4 w-4" />
                        <div class="flex flex-col">
                            <h3 class="text-sm font-medium">Enable SSH</h3>
                            <p class="text-xs text-muted-foreground">
                                Connect to this server via SSH
                            </p>
                        </div>
                    </div>
                    <Switch
                        checked={formData.ssh_enabled}
                        onCheckedChange={(checked) =>
                            (formData.ssh_enabled = checked)}
                    />
                </div>

                {#if formData.ssh_enabled}
                    <div class="space-y-4 pl-2 border-l-2 border-muted mt-4">
                        <div class="grid grid-cols-2 gap-4">
                            <div class="space-y-2">
                                <Label for="username">Username *</Label>
                                <Input
                                    id="username"
                                    bind:value={formData.username}
                                    placeholder="Enter username"
                                />
                            </div>
                            <div class="space-y-2">
                                <Label for="port">Port</Label>
                                <Input
                                    id="port"
                                    type="number"
                                    bind:value={formData.port}
                                    placeholder="22"
                                />
                            </div>
                        </div>

                        <div class="space-y-2">
                            <Label
                                for="security-key"
                                class="flex items-center gap-1"
                            >
                                <KeyRound class="h-3.5 w-3.5" />
                                <span>Security Key *</span>
                            </Label>
                            <div class="relative w-full">
                                <select
                                    id="security-key"
                                    class="w-full h-10 rounded-md border border-input bg-background px-3 py-2 text-sm ring-offset-background focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2 disabled:cursor-not-allowed disabled:opacity-50"
                                    bind:value={formData.security_key}
                                >
                                    <option value="" disabled
                                        >Select security key</option
                                    >
                                    {#if loadingKeys}
                                        <option value="" disabled
                                            >Loading keys...</option
                                        >
                                    {:else if securityKeys.length === 0}
                                        <option value="" disabled
                                            >No security keys available</option
                                        >
                                    {:else}
                                        {#each securityKeys as key}
                                            <option value={key.id}
                                                >{key.name}</option
                                            >
                                        {/each}
                                    {/if}
                                </select>
                            </div>
                            <p class="text-xs text-muted-foreground">
                                Select a security key to use for SSH connection
                            </p>
                        </div>
                    </div>
                {/if}
            </div>

            <div class="flex items-center justify-between pt-2">
                <div class="flex flex-col">
                    <h3 class="text-sm font-medium">Server Active</h3>
                    <p class="text-xs text-muted-foreground">
                        Enable or disable this server
                    </p>
                </div>
                <Switch
                    checked={formData.is_active}
                    onCheckedChange={(checked) =>
                        (formData.is_active = checked)}
                />
            </div>

            <Dialog.Footer class="mt-6">
                <Button type="button" variant="outline" on:click={onClose}>
                    Cancel
                </Button>
                <Button type="submit" disabled={isSubmitting}>
                    {isSubmitting
                        ? selectedServer
                            ? "Updating..."
                            : "Creating..."
                        : selectedServer
                          ? "Update"
                          : "Create"} Server
                </Button>
            </Dialog.Footer>
        </form>
    </Dialog.Content>
</Dialog.Root>
