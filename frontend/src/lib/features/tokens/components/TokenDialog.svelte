<script lang="ts">
    import { Eye, EyeOff } from "lucide-svelte";
    import * as Dialog from "$lib/components/ui/dialog";
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
    import DateTimePicker from "$components/DateTimePicker.svelte";
    import { PROVIDERS, DEFAULT_TOKEN_FORM } from "../constants";
    import type { Token } from "../types";
    import { toast } from "svelte-sonner";

    export let open = false;
    export let onClose: () => void;
    export let onSubmit: (data: any) => Promise<void>; // Changed from handleSubmit to onSubmit
    export let selectedToken: Token | null = null;

    let showTokenValue = false;
    let showRefreshTokenValue = false;
    let formData = selectedToken
        ? { ...selectedToken }
        : { ...DEFAULT_TOKEN_FORM };
    let expiryDate = selectedToken?.expiry
        ? new Date(selectedToken.expiry)
        : null;
    let isSubmitting = false;

    $: if (!open) {
        formData = selectedToken
            ? { ...selectedToken }
            : { ...DEFAULT_TOKEN_FORM };
        expiryDate = selectedToken?.expiry
            ? new Date(selectedToken.expiry)
            : null;
        showTokenValue = false;
        showRefreshTokenValue = false;
    }

    function handleProviderChange(value: string) {
        formData.provider = value;
    }

    function handleExpiryChange(date: Date | null) {
        expiryDate = date;
    }

    async function handleSubmit() {
        console.log("Form Data before submit:", { ...formData, expiryDate }); // Debug log

        if (!formData.provider || !formData.account || !formData.access_token) {
            toast.error("Please fill in all required fields");
            return;
        }

        isSubmitting = true;
        try {
            await onSubmit({
                ...formData,
                expiry: expiryDate?.toISOString() || null,
            });
        } catch (error) {
            console.error("Submit error:", error);
            toast.error("Failed to submit token");
        } finally {
            isSubmitting = false;
        }
    }
</script>

<Dialog.Root bind:open onOpenChange={onClose}>
    <Dialog.Content class="sm:max-w-[425px]">
        <Dialog.Header>
            <Dialog.Title>
                {selectedToken ? "Edit Token" : "Create New Token"}
            </Dialog.Title>
            <Dialog.Description>
                {selectedToken
                    ? "Modify your API token details"
                    : "Add a new API token for integration"}
            </Dialog.Description>
        </Dialog.Header>

        <form class="space-y-4" on:submit|preventDefault={handleSubmit}>
            <div class="space-y-2">
                <Label for="provider">Provider *</Label>
                <Select value={formData.provider}>
                    <SelectTrigger class="w-full">
                        <SelectValue placeholder="Select provider" />
                    </SelectTrigger>
                    <SelectContent>
                        {#each PROVIDERS as provider}
                            <SelectItem
                                value={provider.value}
                                on:click={() =>
                                    handleProviderChange(provider.value)}
                            >
                                {provider.label}
                            </SelectItem>
                        {/each}
                    </SelectContent>
                </Select>
            </div>

            <div class="space-y-2">
                <Label for="account">Account *</Label>
                <Input
                    id="account"
                    bind:value={formData.account}
                    placeholder="Account name or email"
                    required
                />
            </div>

            <div class="space-y-2">
                <Label for="access_token">Access Token *</Label>
                <div class="relative">
                    <Input
                        id="access_token"
                        type={showTokenValue ? "text" : "password"}
                        bind:value={formData.access_token}
                        placeholder="Enter access token"
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

            <div class="space-y-2">
                <Label for="token_type">Token Type</Label>
                <Select bind:value={formData.token_type}>
                    <SelectTrigger class="w-full">
                        <SelectValue placeholder="Select token type" />
                    </SelectTrigger>
                    <SelectContent>
                        <SelectItem value="Bearer">Bearer</SelectItem>
                        <SelectItem value="Basic">Basic</SelectItem>
                        <SelectItem value="OAuth">OAuth</SelectItem>
                    </SelectContent>
                </Select>
            </div>

            <div class="space-y-2">
                <Label for="refresh_token">Refresh Token</Label>
                <div class="relative">
                    <Input
                        id="refresh_token"
                        type={showRefreshTokenValue ? "text" : "password"}
                        bind:value={formData.refresh_token}
                        placeholder="Enter refresh token (optional)"
                    />
                    <button
                        type="button"
                        class="absolute right-2 top-1/2 -translate-y-1/2 text-muted-foreground hover:text-foreground"
                        on:click={() =>
                            (showRefreshTokenValue = !showRefreshTokenValue)}
                    >
                        {#if showRefreshTokenValue}
                            <EyeOff class="w-4 h-4" />
                        {:else}
                            <Eye class="w-4 h-4" />
                        {/if}
                    </button>
                </div>
            </div>

            <div class="space-y-2">
                <Label for="expiry">Expiry Date</Label>
                <DateTimePicker
                    id="expiry"
                    value={expiryDate}
                    on:change={(e) => handleExpiryChange(e.detail)}
                    placeholder="Select expiry date and time"
                />
            </div>

            <div class="space-y-2">
                <Label for="scope">Scope</Label>
                <Input
                    id="scope"
                    bind:value={formData.scope}
                    placeholder="Token scope (optional)"
                />
            </div>

            <Dialog.Footer>
                <Button type="button" variant="outline" on:click={onClose}>
                    Cancel
                </Button>
                <Button type="submit" disabled={isSubmitting}>
                    {isSubmitting
                        ? selectedToken
                            ? "Updating..."
                            : "Creating..."
                        : selectedToken
                          ? "Update"
                          : "Create"} Token
                </Button>
            </Dialog.Footer>
        </form>
    </Dialog.Content>
</Dialog.Root>
