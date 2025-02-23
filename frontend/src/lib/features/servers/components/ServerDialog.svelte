<script lang="ts">
    import { Eye, EyeOff } from "lucide-svelte";
    import * as Dialog from "$lib/components/ui/dialog";
    import { Button } from "$lib/components/ui/button";
    import { Input } from "$lib/components/ui/input";
    import { Label } from "$lib/components/ui/label";
    import * as Select from "$lib/components/ui/select";
    import { PROVIDERS } from "../constants";
    import { DEFAULT_SERVER_FORM } from "../types";
    import type { Server } from "../types";
    import { toast } from "svelte-sonner";

    export let open = false;
    export let onClose: () => void;
    export let onSubmit: (data: any) => Promise<void>;
    export let selectedServer: Server | null = null;

    let showTokenValue = false;
    let formData = selectedServer
        ? { ...selectedServer }
        : { ...DEFAULT_SERVER_FORM };
    let isSubmitting = false;

    $: if (!open) {
        formData = selectedServer
            ? { ...selectedServer }
            : { ...DEFAULT_SERVER_FORM };
        showTokenValue = false;
    }

    function handleProviderChange(value: string) {
        formData.provider = value;
    }

    async function handleSubmit() {
        if (
            !formData.name ||
            !formData.provider ||
            !formData.url ||
            !formData.token
        ) {
            toast.error("Please fill in all required fields");
            return;
        }

        isSubmitting = true;
        try {
            await onSubmit(formData);
        } catch (error) {
            console.error("Submit error:", error);
            toast.error("Failed to submit server");
        } finally {
            isSubmitting = false;
        }
    }
</script>

<Dialog.Root bind:open onOpenChange={onClose}>
    <Dialog.Content class="sm:max-w-[425px]">
        <Dialog.Header>
            <Dialog.Title>
                {selectedServer ? "Edit Server" : "Add New Server"}
            </Dialog.Title>
            <Dialog.Description>
                {selectedServer
                    ? "Modify your server details"
                    : "Add a new server for deployment"}
            </Dialog.Description>
        </Dialog.Header>

        <form class="space-y-4" on:submit|preventDefault={handleSubmit}>
            <div class="space-y-2">
                <Label for="name">Name *</Label>
                <Input
                    id="name"
                    bind:value={formData.name}
                    placeholder="Enter server name"
                    required
                />
            </div>

            <div class="space-y-2">
                <Label for="provider">Provider *</Label>
                <Select.Root
                    value={formData.provider}
                    onValueChange={handleProviderChange}
                >
                    <Select.Trigger class="w-full">
                        <Select.Value placeholder="Select provider" />
                    </Select.Trigger>
                    <Select.Content>
                        {#each PROVIDERS as provider}
                            <Select.Item value={provider.value}>
                                {provider.label}
                            </Select.Item>
                        {/each}
                    </Select.Content>
                </Select.Root>
            </div>

            <div class="space-y-2">
                <Label for="url">URL *</Label>
                <Input
                    id="url"
                    bind:value={formData.url}
                    placeholder="Enter server URL"
                    required
                />
            </div>

            <div class="space-y-2">
                <Label for="token">Token *</Label>
                <div class="relative">
                    <Input
                        id="token"
                        type={showTokenValue ? "text" : "password"}
                        bind:value={formData.token}
                        placeholder="Enter server token"
                        required
                    />
                    <button
                        type="button"
                        class="absolute right-2 top-1/2 -translate-y-1/2 text-muted-foreground hover:text-foreground"
                        on:click={() => (showTokenValue = !showTokenValue)}
                    >
                        {#if showTokenValue}
                            <EyeOff class="w-4 h-4" />
                        {:else}
                            <Eye class="w-4 h-4" />
                        {/if}
                    </button>
                </div>
            </div>

            <Dialog.Footer>
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
